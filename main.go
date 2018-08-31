package main

import (
	"fmt"
	"clmwallet-block-wacther/wacther"
	"sync"
	"os"
	"log"
	"os/signal"
	"syscall"
)

func main()  {

	// 调试strategicpool.go

/*
	n := blocknode.BlockNodeInfo{
	}


	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.First(&n).Error;nil != err {
		log.Println(err)
		log.Println("没有找到")
	}
	log.Println("n.Number:",n.Number)
	log.Println("n.Hash:",n.Hash)
	//
	//var node = &blocknode.BlockNodeInfo{}
	//db := mysqltools.GetInstance().GetMysqlDB()
	//db.First(node, 1) // 查询number为number的node

*/




	/*
	pool := blockpool.Init2()

	n1 := &blocknode.BlockNodeInfo{
		Number:1,
		Hash:"1111111",
	}
	n2 := &blocknode.BlockNodeInfo{
		Number:2,
		Hash:"222222222",
	}
	n3 := &blocknode.BlockNodeInfo{
		Number:3,
		Hash:"333333333",
	}
	n4 := &blocknode.BlockNodeInfo{
		Number:4,
		Hash:"444444444",
	}
	n5 := &blocknode.BlockNodeInfo{
		Number:5,
		Hash:"5555555555",
	}

	pool.InsertElement(n1)
	pool.InsertElement(n2)
	pool.InsertElement(n3)
	pool.InsertElement(n4)
	pool.InsertElement(n5)

	pool.Descrip()

	fmt.Println("-------------")

	pool.RemoveItem(n1)
	pool.ResetEarliestIdx()
	pool.Descrip()

	fmt.Println("-------------")

	pool.RemoveItem(n4)
	pool.ResetEarliestIdx()
	pool.Descrip()

*/

	////////////////////////////////////////////////////////////////////////////


	wacther := wacther.Init()

	var stopLock sync.Mutex
	stop := false
	stopChan := make(chan struct{},1)
	signalChan := make(chan os.Signal,1)
	go func() {
		//阻塞程序运行，直到收到终止的信号
		<- signalChan
		stopLock.Lock()
		stop = true
		stopLock.Unlock()

		log.Println("Cleaning before stop ...")
		wacther.CleanTask()
		stopChan <- struct{}{}
		os.Exit(0)
	}()

	signal.Notify(signalChan,syscall.SIGINT,syscall.SIGTERM,syscall.SIGKILL)



	wacther.WacthStart()




	/////////////////////////////////////////////////////////////////////////////////////////

	/*
			thd := &transactionhandler.EthTransactionHandler{}

			if true == thd.ExistAddress("0xea1771951a97d0fac9378544c4060be4be8b87d3") {
				fmt.Println("找到了")
			} else {
				fmt.Println("没有找到")
			}
	*/

	/*
				if thd.ExistTransByHash("0x4dd0e3092ee7f397ac54975cfe7de374ae1eae124be131b45a9d9803ee03b4e5") {
					fmt.Println("存在")
				} else {
					fmt.Println("不存在")
				}

				err := thd.AddBlockNumberHash("111","xxx","0x4dd0e3092ee7f397ac54975cfe7de374ae1eae124be131b45a9d9803ee03b4e5")
				if nil != err {
					log.Println(err)
				}

				err = thd.InsertReceivedTransInfo("txHash","blockHash1","blockNumber1","fromaddress1","toAddress","0x0011","0x123")
				if nil != err {
					log.Println(err)
				}

				err = thd.InsertReceivedERC20CoinInfo("hxkjk","hashshshh","kjljjsfe","lskjf33","kkkv32","22344","0x32","0x232")
				if nil != err {
					log.Println(err)
				}
			*/

	/*
	resp, err := http.Post("http://47.75.115.210:8781/wallet/Transaction/ExistAddress/",
		"application/x-www-form-urlencoded",
		strings.NewReader("address=0x0001"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
	*/

	//node := &blocknode.BlockNodeInfo{}
	//fmt.Println(node)

	//blocknode.FindAll()

	//wacther := wacther.Init()
	//wacther.WacthStart()
	//wacther.FecthBlockByNumber("0x3bcded")

	//input := "0xa9059cbb0000000000000000000000007c2ee1aaa238bbc3da5d5d74b5dcb1b890b487aa00000000000000000000000000000000000000000000003635c9adc5dea00000"
	//wacther.ERC20ToAddress(input)
	//wacther.ERC20Value(input)




	//var nodes []blocknode.BlockNodeInfo
	////db := mysqltools.GetInstance().GetMysqlDB()
	////db.Find(&nodes)
	//
	//blocknode.Find(&nodes)
	//for _,n := range nodes{
	//	fmt.Println(n)
	//}

	/*
	var node = &blocknode.BlockNodeInfo{}
	db := mysqltools.GetInstance().GetMysqlDB()
	db.First(node, 11) // 查询id为1的product

	fmt.Println(node.Number)
	fmt.Println(node.Hash)
	*/



	fmt.Println("End")

	/*
	wg := &sync.WaitGroup{}
	wg.Add(1)

	var c chan string = make(chan string,100)
	go func() {
		for {
			for i := 1; i < 10; i++ {
				c <- "abc" + strconv.Itoa(i)
			}
			time.Sleep(2 * time.Second)
		}
		//wg.Done()
	}()

	go func() {
		for   {
			for v:= range c {
				fmt.Println(v)
			}
		}
	}()

	wg.Wait()
	*/
	/*
	fmt.Println("Block Wacther is Running...")
	var wthr = wacther2.BlockWacther{}
	//wthr.StartWacth()
	wthr.Dooo = func(i int) int {
		fmt.Println("I'm a blocker func")
		return i
	}
	wthr.Dooo(1)
	fmt.Println("End")
	*/

	/*
	str := "0xf"
	intStr,_ := strconv.ParseInt(str,0,64)
	fmt.Println(intStr)


	str = strconv.FormatInt(intStr,3)
	fmt.Println(str)
	*/

}
