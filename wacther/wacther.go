package wacther

import (
	"fmt"
	//"sync"
	"time"
	"clmwallet-block-wacther/blockpool"
	"clmwallet-block-wacther/modles/blocknode"
	"github.com/ethereum/go-ethereum/rpc"
	"clmwallet-block-wacther/clinterface"
	"clmwallet-block-wacther/config"
	"strconv"
	"reflect"
)




type BlockWacther struct {
	blockPool *blockpool.BlockPool

	affirmChain chan *blocknode.BlockNodeInfo
	resendChain chan *blocknode.BlockNodeInfo

	client *rpc.Client

	TransHandler clinterface.TransInterface

}

func Init() *BlockWacther {
	b := &BlockWacther{
		blockPool:blockpool.InitBlockPool(),
		affirmChain:make(chan *blocknode.BlockNodeInfo,1000),
		resendChain:make(chan *blocknode.BlockNodeInfo,1000),
	}
	return b
}


//var wg sync.WaitGroup

func (bw BlockWacther) StartWacth()  {
	//wg.Add(1)
	//go bw.Fecther2Parse()
	//go bw.AffirmTranscations()
	//go bw.ReSendTranscations()
	//wg.Wait()
}





func (bw BlockWacther)FecthAndParseBlock(currentBlockNumber int64)  {

	//获取最新的区块
	ethLastNode := bw.FecthBlockByNumber("latest")


	//获取eth最新区块与本地同步的最新区块的高度之差
	var diff int64 = 0
	if !bw.blockPool.IsEmpty() {
		diff = ethLastNode.BlockNumber - bw.blockPool.LatestNumber() + config.CLAffiremBlockHeigh
	}

	//由当前区块向前拉取区块
	var count int64 = 0
	var nodeHash = ethLastNode.BlockHash
	for ; count < diff;count++  {
		node := bw.FecthBlockByHash(nodeHash)
		bw.putBlock2affirmChain(node)
		bw.putBlock2resendChain(node)

		nodeHash = node.ParentHash
	}

}

// 选出需要重新打包发送的区块号
func (bw BlockWacther) putBlock2resendChain(info *blocknode.BlockNodeInfo)  {
	if nil == bw.resendChain {
		return
	}

	node := bw.blockPool.ReciveBlock(info)
	if node != nil {
		bw.resendChain <- node
	}

}

//选出需要确认的
func (bw BlockWacther) putBlock2affirmChain(info *blocknode.BlockNodeInfo)  {
	if nil == bw.affirmChain {
		return
	}

	ok,node := bw.blockPool.LookBocks4AffirmTrans()
	if !ok {
		return
	}

	if node != nil {
		bw.affirmChain <- node
	}

}
// 通过区块号拉取区块
func (bw BlockWacther) FecthBlockByNumber(blockNumber string) *blocknode.BlockNodeInfo {

	client,err := bw.getClient()
	if nil != err {
		return nil
	}

	var blockInfo = make(map[string]interface{})
	if err := client.Call(&blockInfo,"eth_getBlockByNumber",blockNumber,"true");err != nil {
		return nil
	}

	bw.parseBlock(blockInfo) //解析区块

	number,_ := blockInfo["blockNumber"].(string)
	numberInt,_ := strconv.ParseInt(number,0,64)
	hash,_ := blockInfo["blockHash"].(string)
	parentHash,_ := blockInfo["parentHash"].(string)

	node := &blocknode.BlockNodeInfo{
		BlockNumber:numberInt,
		BlockHash:hash,
		ParentHash:parentHash,
	}

	return node
}

// 通过区块HASH拉取区块
func (bw BlockWacther) FecthBlockByHash(blockHash string) *blocknode.BlockNodeInfo{

	client,err := bw.getClient()
	if nil != err {
		return nil
	}

	var blockInfo = make(map[string]interface{})
	if err := client.Call(&blockInfo,"eth_getBlockByHash",blockHash,"true");err != nil {
		return nil
	}

	bw.parseBlock(blockInfo)

	number,_ := blockInfo["blockNumber"].(string)
	numberInt,_ := strconv.ParseInt(number,0,64)
	hash,_ := blockInfo["blockHash"].(string)
	parentHash,_ := blockInfo["parentHash"].(string)

	node := &blocknode.BlockNodeInfo{
		BlockNumber:numberInt,
		BlockHash:hash,
		ParentHash:parentHash,
	}

	return node

}

///解析区块信息
func (bw BlockWacther)parseBlock(blockInfo map[string]interface{})  {
	if nil == blockInfo { return }

	//得到交易信息数组
	transInfoI := blockInfo["transactions"]
	if nil == transInfoI {return }
	fmt.Println(reflect.TypeOf(transInfoI))
	transInfo,ok := transInfoI.([]interface{})
	if !ok {
		return
	}

	fmt.Println(transInfo)

	//	解析每一个交易
	// 	WARNING:伪代码,具体字段需要确认
	for _,mI := range transInfo {
		fmt.Println(reflect.TypeOf(mI))
		m,ok := mI.(map[string]interface {})
		if !ok {
			break
		}

		blockHash := m["blockHash"].(string)
		fmt.Println(blockHash)

		blockNumber := m["blockNumber"].(string)
		fmt.Println(blockNumber)

		//transactionIndex := m["transactionIndex"].(string)

		hash := m["hash"].(string)
		fmt.Println(hash)

		//nonce := m["nonce"].(string)

		from := m["from"].(string)
		fmt.Println(from)

		to := m["to"].(string)
		fmt.Println(to)

		value := m["value"].(string)
		fmt.Println(value)

		//gas := m["gas"].(string)
		//gasPrice := m["gasPrice"].(string)
		//input := m["input"].(string)


		/*
		//根据地址判断是不是属于超链平台上的用户
		if nil != bw.TransHandler {
			if bw.TransHandler.ExistAddress(from) { //是我们发出的交易
				bw.TransHandler.AddBlockNumberHash(blockNumber,blockHash,tHash)
			}

			if bw.TransHandler.ExistAddress(to) {//是我们平台发出去的交易
				//根据交易Hash 增加blockNumber/blockHash到交易数据库表
				bw.TransHandler.InsertTransInfo(from,to,value,gas,tHash)
			}
		}
		*/
	}
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







// 确认交易
func (bw BlockWacther) AffirmTranscations()  {
	//通过chain实现休眠
	if nil == bw.affirmChain {
		return
	}
	for {
		for v := range bw.affirmChain {
			fmt.Println(v)
		}
	}

}

// 重发交易 仅通知 相关的API重发，不作具体重新发送功能
func (bw BlockWacther)ReSendTranscations()  {
	// 访问交易表
	//通过chain实现休眠
	for {
		fmt.Println("ReSend Transcations ...")
		time.Sleep(time.Second * 5)
	}
}
