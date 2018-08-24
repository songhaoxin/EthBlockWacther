package transactionhandler

import (
	"github.com/kirinlabs/HttpRequest"
	"log"
	"encoding/json"
)

type EthTransactionHandler struct {

}

// 判断指定的 地址 是否属于平台内帐户
func (t *EthTransactionHandler) ExistAddress(address string) bool {

	req := HttpRequest.NewRequest()

	res,err := req.Post("http://47.75.115.210:8781/wallet/Transaction/ExistAddress",map[string]interface{} {
		"address":address,
	})
	if nil != err {
		log.Println(err)
		return false
	}


	if 200 != res.StatusCode() {
		return false
	}

	body,err := res.Body()
	if nil != err {
		log.Println(err)
		return false
	}


	resMap := make(map[string]interface{})
	err = json.Unmarshal(body,&resMap)
	if nil != err {
		log.Println(err)
		return false
	}
	if v,ok := resMap["status"].(float64);ok {
		if 0 == v {
			return true
		}
	}


	return false

}