package main

import (
	"fmt"
	"clmwallet-block-wacther/wacther"
)

func main()  {

	//thd := &transactionhandler.EthTransactionHandler{}
	//if true == thd.ExistAddress("0xbd88ba7a14b5263953f236d0182c311ff1102829") {
	//	fmt.Println("打到了")
	//} else {
	//	fmt.Println("没有找到")
	//}


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

	wacther := wacther.Init()
	wacther.WacthStart()


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
