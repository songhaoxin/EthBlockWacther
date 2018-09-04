/*
@Time : 2018/9/4 上午11:00 
@Author : Mingjian Song
@File : helper
@Software: 深圳超链科技
*/

package helper

import (
	"strings"
	"strconv"
	"log"
	"math"
	"github.com/shopspring/decimal"
	"math/big"
	"errors"
)

//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
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
func Substr2(str string, start int, end int) string {
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

// Hexadecimal to decimal
func HexDec(h string) (n int64) {
	s := strings.Split(strings.ToUpper(h), "")
	l := len(s)
	i := 0
	d := float64(0)
	hex := map[string]string{"A": "10", "B": "11", "C": "12", "D": "13", "E": "14", "F": "15"}
	for i = 0; i < l; i++ {
		c := s[i]
		if v, ok := hex[c]; ok {
			c = v
		}
		f, err := strconv.ParseFloat(c, 10)
		if err != nil {
			log.Println("Hexadecimal to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(16, float64(l-i-1))
	}
	return int64(d)
}

func Hex2Decimal(hex string,tokenDecimal int,round int)  (decimal.Decimal,error){

	sv := string([]byte(hex)[2:])
	bv,ok := new(big.Int).SetString(sv,16)
	if !ok {
		return decimal.NewFromBigInt(big.NewInt(0),0),errors.New("Hex string convert to big.Int error!")
	}
	dv := decimal.NewFromBigInt(bv,0)

	decimalV := decimal.NewFromBigInt(big.NewInt(1),int32(tokenDecimal))

	dv = dv.Div(decimalV).Round(2)

	return dv,nil
}
