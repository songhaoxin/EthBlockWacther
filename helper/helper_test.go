/*
@Time : 2018/9/4 下午12:32 
@Author : Mingjian Song
@File : helper_test
@Software: 深圳超链科技
*/

package helper

import (
	"testing"
	"fmt"
)

func TestHex2Decimal(t *testing.T) {
	dv,err := Hex2Decimal("0xAAA56bc75e2d6310000",18,6)
	fmt.Println(dv,err)
}
