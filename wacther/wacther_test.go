/*
@Time : 2018/9/5 下午10:44 
@Author : Mingjian Song
@File : wacther_test
@Software: 深圳超链科技
*/

package wacther

import (
	"testing"


)

func TestFecthBlockByNumber(t *testing.T)  {
	wacther := Init()
	wacther.FecthBlockByNumber("0x3cc507")
}
