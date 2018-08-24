package wacther

import (

	"clmwallet-block-wacther/blockpool"
	"clmwallet-block-wacther/modles/blocknode"
	"github.com/ethereum/go-ethereum/rpc"
	"clmwallet-block-wacther/clinterface"
	"clmwallet-block-wacther/config"
	"strconv"
	"strings"
	"sync"
	"time"
	"log"
)




type BlockWacther struct {
	blockPool *blockpool.BlockPool

	affirmChain chan string
	resendChain chan string

	client *rpc.Client
	fecthTimer *time.Timer

	TransHandler clinterface.TransInterface

}

func Init() *BlockWacther {
	b := &BlockWacther{
		blockPool:blockpool.Init(),
		affirmChain:make(chan string,1000),
		resendChain:make(chan string,1000),
		fecthTimer:time.NewTimer(config.TimeDelayInSecand * time.Second),
	}
	return b
}


var wg sync.WaitGroup

func (bw BlockWacther) WacthStart()  {
	wg.Add(1)
	go bw.fecthParseBlock()
	go bw.AffirmTranscations()
	go bw.ReSendTranscations()
	wg.Wait()
}

// 确认交易
func (bw BlockWacther) AffirmTranscations()  {
	if nil == bw.affirmChain {
		return
	}
	for {
		for th := range bw.affirmChain {
			log.Println("处理已经确认的交易:",th)
			if  nil != bw.TransHandler {
				bw.TransHandler.AffirmTrans(th)
			}
		}
	}

}

// 重发交易
func (bw BlockWacther)ReSendTranscations()  {
	if nil == bw.resendChain {return }
	for {
		for th := range bw.resendChain {
			log.Println("处理需要重新打包的交易:",th)
			if nil != bw.TransHandler {
				bw.TransHandler.ResendTrans(th)
			}
		}
	}
}



/// 从geth结点上同步最新区块，并解析区块，实现交易确认、处理需要重新打包的交易
func (bw BlockWacther)fecthParseBlock()  {
	for {
		select {
		case <-bw.fecthTimer.C:
			log.Println("开始同步区块")
			bw.fecthTimer.Stop()
			bw._fecthParseBlock()

			//bw.affirmChain <- "affirm!"
			//bw.resendChain <- "resend!"
			bw.fecthTimer.Reset(time.Second * config.TimeDelayInSecand)
		}
	}
}

func (bw BlockWacther) _fecthParseBlock() {
	//获取最新的区块
	ethLastNode := bw.fecthBlockByNumber("latest")

	if bw.blockPool.ContainElement(ethLastNode) {
		log.Printf("当前没有更新的区块")
		return

	}



	var needFecthCount int64 = 0
	if !bw.blockPool.IsEmpty() {
		needFecthCount = ethLastNode.Number - bw.blockPool.LatestNumber() - 1
	}
	needFecthCount += config.AffiremBlockHeigh

	var fecthCount int64 = 0
	node := ethLastNode
	hash := node.ParentHash
	for ; fecthCount < needFecthCount;fecthCount++  {
		bw.putBlock2resendChain(node)
		node = bw.fecthBlockByHash(hash)
		hash = node.ParentHash
	}

	// 选出已经确认的交易
	bw.putBlock2affirmChain()
}

// 选出需要重新打包发送的区块号
func (bw BlockWacther) putBlock2resendChain(info *blocknode.BlockNodeInfo)  {

	if nil == info {return }

	node := bw.blockPool.ReciveBlock(info)
	if nil == node || "" == node.TransHash {return }

	tHashs := strings.Split(node.TransHash, ";")
	for _,tH := range tHashs {
		if "" != tH {
			if nil != bw.resendChain {
				bw.resendChain <- tH
			}
		}
	}
}

//选出需要确认的交易
func (bw BlockWacther) putBlock2affirmChain()  {

	txHashs := bw.blockPool.LookBocks4AffirmTrans()

	for _,txHash := range txHashs {
		if "" != txHash {
			if nil != bw.affirmChain {
				bw.affirmChain <- txHash
			}
		}
	}
}


// 通过区块号拉取区块
func (bw BlockWacther) fecthBlockByNumber(blockNumber string) *blocknode.BlockNodeInfo {

	client,err := bw.getClient()
	if nil != err {
		return nil
	}

	var blockInfo = make(map[string]interface{})
	if err := client.Call(&blockInfo,"eth_getBlockByNumber",blockNumber,"true");err != nil {
		return nil
	}

	number,_ := blockInfo["number"].(string)
	numberInt,_ := strconv.ParseInt(number,0,64)
	hash,_ := blockInfo["hash"].(string)
	parentHash,_ := blockInfo["parentHash"].(string)

	node := &blocknode.BlockNodeInfo{
		Number:numberInt,
		Hash:hash,
		ParentHash:parentHash,
	}

	transHashs := ""
	if false == bw.blockPool.ContainElement(node) {
		transHashs = bw.parseBlock(blockInfo) //解析区块
	}
	node.TransHash = transHashs

	return node
}

// 通过区块HASH拉取区块
func (bw BlockWacther) fecthBlockByHash(blockHash string) *blocknode.BlockNodeInfo{

	client,err := bw.getClient()
	if nil != err {
		return nil
	}

	var blockInfo = make(map[string]interface{})
	if err := client.Call(&blockInfo,"eth_getBlockByHash",blockHash,"true");err != nil {
		return nil
	}

	number,_ := blockInfo["number"].(string)
	numberInt,_ := strconv.ParseInt(number,0,64)
	hash,_ := blockInfo["hash"].(string)
	parentHash,_ := blockInfo["parentHash"].(string)

	node := &blocknode.BlockNodeInfo{
		Number:numberInt,
		Hash:hash,
		ParentHash:parentHash,
	}

	transHashs := ""
	if false == bw.blockPool.ContainElement(node) {
		transHashs = bw.parseBlock(blockInfo) //解析区块
	}
	node.TransHash = transHashs

	return node

}

///解析区块信息
func (bw BlockWacther)parseBlock(blockInfo map[string]interface{}) string {
	if nil == blockInfo { return "" }

	//得到交易信息数组
	transInfoI := blockInfo["transactions"]
	if nil == transInfoI {return ""}

	transInfo,ok := transInfoI.([]interface{})
	if !ok {
		return ""
	}


	// 保存本次区块中所包含的与本平台帐户相关的 '交易hash'
	var transHashs string = ""

	//	解析每一个交易
	for _,mI := range transInfo {
		m,ok := mI.(map[string]interface {})
		if !ok {
			break
		}

		blockHash := m["blockHash"].(string)
		blockNumber := m["blockNumber"].(string)
		//transactionIndex := m["transactionIndex"].(string)
		hash := m["hash"].(string)
		//nonce := m["nonce"].(string)
		from := m["from"].(string)
		to := m["to"].(string)
		value := m["value"].(string)
		//gas := m["gas"].(string)
		//gasPrice := m["gasPrice"].(string)
		//input := m["input"].(string)


		//根据地址判断是不是属于超链平台上的用户
		if nil != bw.TransHandler {
			if bw.TransHandler.ExistAddress(from) { //是我们发出的交易
				// 对已经发出去的交易，填充好区块号及区块Hash
				bw.TransHandler.AddBlockNumberHash(blockNumber,blockHash,hash)
			}

			if bw.TransHandler.ExistAddress(to) {//别人向本平台帐户转帐
				//根据交易Hash 增加blockNumber/blockHash到交易数据库表
				bw.TransHandler.InsertReceiveTransInfo(hash,blockHash,blockNumber,from,to,value)
			}

			transHashs = transHashs + hash + ";"
		}

	}

	return transHashs
}


func (bw BlockWacther) getClient() (client *rpc.Client,err error) {
	err = nil
	if nil == bw.client{
		bw.client,err = rpc.Dial(config.GethHost)
	}
	return bw.client,err
}



//获取本地同步最新的区块
func (bw BlockWacther) getNewestBlockNumber() int64  {
	return bw.blockPool.LatestNumber()
}


