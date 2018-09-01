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



	fmt.Println("End")

}
