package main

import "fmt"

func main()  {

	//node := &blocknode.BlockNodeInfo{}
	//fmt.Println(node)
	type N struct {
		A string
		B string
		C string
	}

	m := make(map[int64] *N)
	n := &N{}
	n.A = "A"
	n.B = "B"

	m[1234] = n


	n1 := m[1234]

	fmt.Println(n1.A)
	fmt.Println(n1.B)

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
