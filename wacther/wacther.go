package wacther

import (

	"clmwallet-block-wacther/blockpool"
	"clmwallet-block-wacther/modles/blocknode"
	"github.com/ethereum/go-ethereum/rpc"
	"clmwallet-block-wacther/clinterface"
	"clmwallet-block-wacther/configs"
	"strconv"
	"strings"
	"sync"
	"time"
	"log"
	"clmwallet-block-wacther/transactionhandler"
)




type BlockWacther struct {

	// 维持同步的区块信息的存储池
	// 该池中保存着足以用来确定交易的区块链信息
	//blockPool *blockpool.BlockPool
	blockPool *blockpool.StrategicPool


	// 用来传递 需要确定 的交易hash 的通道
	// 由专门的服务去从该通道中读取交易做进一步的处理
	succeedBlocksChain chan string

	// 用来传递 需要重新发送 的交易hash 的通道
	// 由专门的服务去从该通道中读取交易做进一步的处理
	failedBlocksChain chan string

	// 用来从geth结点中同步区块的网络客户实例
	client *rpc.Client

	// 定时器，实现定时从geth结点中拉取区块信息
	fecthTimer *time.Timer

	// 满足TransInterface接口的交易处理器
	// 负责具体处理 确认交易、重发交易 等功能
	// 之所以设计成 接口方式，是为了解耦 拉取区块功能 和 处理交易功能
	// 这样，该模块可以轻松的重用到 比特币 等其他的平台上去
	TransHandler clinterface.TransInterface

}

func Init() *BlockWacther {
	b := &BlockWacther{
		blockPool:          blockpool.Init(),
		succeedBlocksChain: make(chan string,1000),
		failedBlocksChain:  make(chan string,1000),
		fecthTimer:         time.NewTimer(configs.TimeDelayInSecand * time.Second),
	}

	b.TransHandler = transactionhandler.Init()

	return b
}



var wg sync.WaitGroup

/// 开启监视程序
func (bw BlockWacther) WacthStart()  {
	// 不调用它的.Done()方法，以另类地实现循环
	wg.Add(1)
	go bw.FecthParseBlockFromServerTimes()
	go bw.HandleSuccessedTrans()
	go bw.HandleFailedTrans()
	wg.Wait()
}

/// 在非硬件原因宕机的情况，确保程序在退出之前执行清理工作
func (bw BlockWacther) CleanTask() {
	// 把确认的成功交易数据全部通知客户端
	// 把失败的交易全部通知到客户端
	//bw.blockPool.Save2Db()
}


// 确认交易
func (bw BlockWacther) HandleSuccessedTrans()  {
	if nil == bw.succeedBlocksChain {
		return
	}
	for {
		for th := range bw.succeedBlocksChain {
			if  nil != bw.TransHandler {
				log.Println("发了成功短信的交易HASH:",th)
				bw.TransHandler.NoticeTransAffirmed(th)
			}
		}
	}
}

// 重发交易
func (bw BlockWacther) HandleFailedTrans()  {
	if nil == bw.failedBlocksChain {return }
	for {
		for th := range bw.failedBlocksChain {

			if nil != bw.TransHandler {
				log.Println("发了失败短信的交易HASH:",th)
				bw.TransHandler.NoticeTransFailed(th)
			}
		}
	}
}



/// 从geth结点上同步最新区块，并解析区块，实现交易确认、处理需要重新打包的交易
func (bw BlockWacther) FecthParseBlockFromServerTimes()  {
	for {
		select {
		case <-bw.fecthTimer.C:
			//log.Println("开始同步区块")
			bw.fecthTimer.Stop()
			bw._fecthParseBlock()

			//bw.succeedBlocksChain <- "affirm!"
			//bw.failedBlocksChain <- "resend!"
			bw.fecthTimer.Reset(time.Second * configs.TimeDelayInSecand)
		}
	}
}

//从服务端获取的 需要确认的最早的区块号
func (bw BlockWacther) getStartIdxFromServer() int64{
	return 0
}

func (bw BlockWacther) _fecthParseBlock() {
	//获取最新的区块
	ethLastNode := bw.FecthBlockByNumber("latest")
	
	if nil == ethLastNode || ethLastNode.Number < 0 || "" == ethLastNode.Hash {
		return 
	}

	if bw.blockPool.ContainElement(ethLastNode) {
		//log.Printf("当前没有更新的区块")
		return
	}

	// 更新最新的区块号
	bw.blockPool.SetLatestIdx(ethLastNode.Number)
	//筛掉失败的区块交易
	bw.ReceiveBlocksFilterFailed(ethLastNode)


	var needFecthCount int64 = 0
	if 0 < bw.blockPool.Size() {
		needFecthCount = ethLastNode.Number - bw.blockPool.GetEarliestIdx()
	}

	// 向前回溯获取区块
	var fecthCount int64 = 0
	node := ethLastNode
	hash := node.ParentHash
	for ; fecthCount < needFecthCount;fecthCount++  {

		hashValueNuber,_ :=  strconv.ParseInt(hash,0,64)
		if hashValueNuber <= 0 {break}

		node = bw.fecthBlockByHash(hash)
		if nil == node {break}
		hash = node.ParentHash
		//筛掉失败的区块交易
		bw.ReceiveBlocksFilterFailed(node)
	}

	// 选出已经确认的交易
	bw.HandleBlocksSuccessed()
}

// 选出需要重新打包发送的区块号
func (bw BlockWacther) ReceiveBlocksFilterFailed(info *blocknode.BlockNodeInfo)  {

	if nil == info {return }

	node := bw.blockPool.ReciveBlockFromChain(info)
	if nil == node {return }



	if "" == node.TransHash {return }

	tHashs := strings.Split(node.TransHash, ";")

	for _,tH := range tHashs {
		if "" != tH {
			if nil != bw.failedBlocksChain {
				bw.failedBlocksChain <- tH
				log.Println("往chain中发送失败的交易,区块号:",node.Number)
				log.Println("往chain中发送失败的交易,区块HASH:",tH)
			}
		}
	}
}

//选出需要确认的交易
func (bw BlockWacther) HandleBlocksSuccessed()  {

	txHashs := bw.blockPool.LookSuccessedTransHashs()

	for _,txHash := range txHashs {
		if "" != txHash {
			if nil != bw.succeedBlocksChain {
				bw.succeedBlocksChain <- txHash
				log.Println("发送成功交易HASH:%s",txHashs)

			}
		}
	}
}


// 通过区块号拉取区块
func (bw BlockWacther) FecthBlockByNumber(blockNumber string) *blocknode.BlockNodeInfo {

	client,err := bw.getClient()
	if nil != err {
		return nil
	}

	var blockInfo = make(map[string]interface{})
	if err := client.Call(&blockInfo,"eth_getBlockByNumber",blockNumber,true);err != nil {
		log.Println(err)
		return nil
	}

	// 如果没有找到区块信息
	if 0 == len(blockInfo) {
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
		log.Println(err)
		return nil
	}

	var blockInfo = make(map[string]interface{})
	if err := client.Call(&blockInfo,"eth_getBlockByHash",blockHash,true);err != nil {
		log.Println(err)
		return nil
	}

	// 如果没有找到区块信息
	if 0 == len(blockInfo) {
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


func (bw BlockWacther) getClient() (client *rpc.Client,err error) {
	err = nil
	if nil == bw.client{
		bw.client,err = rpc.Dial(configs.GethHost)
	}
	return bw.client,err
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
			gas := m["gas"].(string)
			//gasPrice := m["gasPrice"].(string)
			input := m["input"].(string)


			// 如果没有 交易处理 的代理直接返回
			if nil == bw.TransHandler {
				return ""
			}

			// 如果既不是 以太币 也不是 ERC20 交易
			if (!isEthTransf(value)) && (!isERC20Transf(value,input)) {
				continue
			}

			// 先判断发出的交易，对于发出的交易的，只按发出的地址进行确认
			if bw.TransHandler.ExistAddress(from) { //是我们发出的交易
				log.Println("是我们平台发出的交易：",hash)
				// 对已经发出去的交易，填充好区块号及区块Hash
				transHashs = transHashs + hash + ";"
				if nil == bw.TransHandler.AddBlockNumberHash(blockNumber,blockHash,hash) {
					log.Println("填充交易的blockNumber 和 blockHash 成功")
				} else {
					log.Println("填充交易的blockNumber 和 blockHash 失败")
				}

				// 对于我们发出去的交易，只以保存一条记录到数据表中，所以直接解析下一条
				continue
			}



			//再判断是不是别人发给我们的交易
			//如果是ERC20代币交易
			if isERC20Transf(value,input) {
				erc20to := erc20ToAddress(input)
				erc20Value := erc20Value(input)

				//根据交易Hash 增加blockNumber/blockHash到交易数据库表
				if bw.TransHandler.ExistAddress(erc20to) { //给本平台帐户发送ERC20代币
					log.Println("收到别人发给我们代币交易：",hash)
					transHashs = transHashs + hash + ";"
					if nil == bw.TransHandler.InsertReceivedERC20CoinInfo(hash,blockHash,blockNumber,from,erc20to,to,gas,erc20Value) {
						log.Println("保存 接收 代币 交易信息 成功")
					}
				}
			} else if isEthTransf(value) {
				//根据交易Hash 增加blockNumber/blockHash到交易数据库表
				if bw.TransHandler.ExistAddress(to) {
					log.Println("收到别人发给我们的以太币交易：",hash)
					transHashs = transHashs + hash + ";"
					if nil == bw.TransHandler.InsertReceivedTransInfo(hash,blockHash,blockNumber,from,to,gas,value) {
						log.Println("保存 接收 以太币 交易信息 成功")
					}
				}
			}

	}


	return transHashs
}




////获取本地同步最新的区块
//func (bw BlockWacther) getNewestBlockNumber() int64  {
//	return bw.blockPool.LatestNumber()
//}

// 是否是ERC20代币交易
func isERC20Transf(value string,input string)  bool  {
	return (value == "0x0")  && (substr(input,0,10) == "0xa9059cbb")
}

// 是否是以太币交易
func isEthTransf(value string) bool {
	valueNum,err :=  strconv.ParseInt(value,0,64)
	if nil != err {
		return false
	}
	if valueNum <= 0 {
		return false
	}

	return true
}

// 获取ERC20代币的发送地址
func erc20ToAddress(input string) string  {
	if 138 != len(input) {
		return ""
	}
	address := "0x" + substr(input,34,40)
	return address
}

// 获取ERC20代币的发送值
func erc20Value(input string) string  {
	if 138 != len(input) {
		return ""
	}
	value := substr2(input,74,137)
	return value
}

//截取字符串 start 起点下标 length 需要截取的长度
func substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func substr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}


