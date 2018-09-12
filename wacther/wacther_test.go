/*
@Time : 2018/9/5 下午10:44 
@Author : Mingjian Song
@File : wacther_test
@Software: 深圳超链科技
*/

package wacther

import (
	"testing"

	"fmt"
)

func TestFecthBlockByNumber(t *testing.T)  {
	wacther := Init()
	wacther.FecthBlockByNumber("0x3cc507")
}

func TestBlockWacther_CleanTask(t *testing.T) {
	value := erc20Value("0xa9059cbb00000000000000000000000056dbeaf802dc6ab440d99bb4e866596aa77b468a0000000000000000000000000000000000000000000000056bc75e2d63100000")
	fmt.Println(value)
}