/*
@Time : 2018/9/1 上午11:38 
@Author : Mingjian Song
@File : strategicpool_test
@Software: 深圳超链科技
*/

package blockpool

import (
	"testing"
	"database/sql"
	"fmt"
)


func TestResetEarliestIdx(t *testing.T)  {

	db, err := sql.Open("mysql","root:root@tcp(120.77.223.246:3306)/clwallet")
	if nil != err {
		return
	}

	rows,err := db.Query("SELECT 1 FROM blockNodeInfo WHERE number=60 LIMIT 1")
	if nil != err {
		return
	}

	for rows.Next() {
		columns,_ := rows.Columns()
		fmt.Println(columns)
	}

	db.Close()

}
