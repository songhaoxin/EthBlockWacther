/*
@Time : 2018/9/4 下午12:32 
@Author : Mingjian Song
@File : helper_test
@Software: 深圳超链科技
*/

package helper

import (
	"testing"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"log"
)

func TestHex2Decimal(t *testing.T) {
	d,e := Hex2Decimal("0x38d7ea4c68000",18,8)
	if nil != e {
		log.Println(e)
		fmt.Println(d.String())
	}
	fmt.Println(d.String())
	//从服务端获取的 需要确认的最早的区块号

	/*
		db,err := sql.Open("mysql",configs.ServerDBConnectString)
		if nil != err {
			log.Println("err")
			return
		}
		defer db.Close()

		var lowestNumber int64
		err = db.QueryRow(`SELECT block_number FROM send_out_detail  WHERE is_confirm = 0 ORDER BY block_number ASC LIMIT 0,1`).Scan(&lowestNumber)
		if nil != err {

		}

		log.Println(lowestNumber)
*/

}
