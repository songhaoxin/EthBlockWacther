package main

import (
	"strconv"
	"fmt"
)

func main()  {
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

	str := "0xf"
	intStr,_ := strconv.ParseInt(str,0,64)
	fmt.Println(intStr)


	str = strconv.FormatInt(intStr,3)
	fmt.Println(str)


}
