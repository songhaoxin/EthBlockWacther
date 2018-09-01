package transactionhandler

import (
	"github.com/kirinlabs/HttpRequest"
	"log"
	"encoding/json"
	"clmwallet-block-wacther/configs"
	"errors"
)

type EthTransactionHandler struct {

}

func Init() *EthTransactionHandler {
	return &EthTransactionHandler{}
}

const (
	ExistAddressAPI                = configs.ServerHost + "/wallet/Transaction/ExistAddress"
	ExistTransByHashAPI            = configs.ServerHost + "/wallet/Transaction/ExistTransByHash"
	AddBlockNumberHashAPI          = configs.ServerHost + "/wallet/Transaction/AddBlockNumberHash"
	InsertReceivedTransInfoAPI     = configs.ServerHost + "/wallet/Transaction/InsertReceiveTransInfo"
	InsertReceivedERC20CoinInfoAPI = configs.ServerHost + "/wallet/Transaction/InsertReceivedERC20CoinInfo"
	NoticeTransSucceededAPI        = configs.ServerHost + "/wallet/Transaction/NoticeTransAffirmed"
	NoticeTransFailedAPI           = configs.ServerHost + "/wallet/Transaction/NoticeTransFailed"
	GetUnHandledTransInfoAPI       = configs.ServerHost + "/wallet/Transaction/GetUnHandledTransInfo"
)

//var req = HttpRequest.NewRequest()

// 判断指定的 地址 是否属于平台内帐户
func (t EthTransactionHandler) ExistAddress(address string) bool {

	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})

	res,err := req.Post(ExistAddressAPI,map[string]interface{} {
		"address":address,
	})

	if nil != err {
		log.Println(err)
		return false
	}


	log.Println(res.StatusCode())
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

func (t EthTransactionHandler) ExistTransByHash(transHash string) bool {
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})

	res,err := req.Post(ExistTransByHashAPI,map[string]interface{} {
		"txHash":transHash,
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

//根据交易Hash 增加blockNumber/blockHash到交易数据库表
func (t EthTransactionHandler) AddBlockNumberHash(blockNumber string,blockHash string,withTransHash string) error{

	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})

	res,err := req.Post(AddBlockNumberHashAPI,map[string]interface{} {
		"blockNumber":blockNumber,
		"blockHash":blockHash,
		"withTransHash":withTransHash,
	})

	if nil != err {
		return err
	}

	body,err := res.Body()
	if nil != err {
		return err
	}


	resMap := make(map[string]interface{})
	err = json.Unmarshal(body,&resMap)
	if nil != err {
		return err
	}

	stuCode,_  := resMap["status"].(float64)
	msg,_ := resMap["msg"].(string)

	if 0 == stuCode && 200 == res.StatusCode() {
		return nil
	}


	return errors.New(msg)

}

//根据解析的交易的信息（别人向我们的帐户转帐这一情况），添加一条记录到交易表
func (t EthTransactionHandler)InsertReceivedTransInfo(
	hash string,
	blockHash string,
	blockNumber string,
	fromAddress string,
	toAddress string,
	gas string,
	value string) error {

	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})

	res,err := req.Post(InsertReceivedTransInfoAPI,map[string]interface{} {
		"txHash":hash,
		"blockHash":blockHash,
		"blockNumber":blockNumber,
		"fromAddress":fromAddress,
		"toAddress":toAddress,
		"value":value,
		"types":"ETH",
		"gas":gas,
	})

	if nil != err {
		return err
	}

	body,err := res.Body()
	if nil != err {
		return err
	}

	resMap := make(map[string]interface{})
	err = json.Unmarshal(body,&resMap)
	if nil != err {
		return err
	}


	stuCode,_  := resMap["status"].(float64)
	msg,_ := resMap["msg"].(string)

	if 0 == stuCode && 200 == res.StatusCode() {
		return nil
	}

	return errors.New(msg)

}

//根据解析的ERC20代币交易信息(别人向我们的帐户转帐这一情况），添加一条记录到交易表
func (t EthTransactionHandler)  InsertReceivedERC20CoinInfo(
	hash string,
	blockHash string,
	blockNumber string,
	fromAddress string,
	toAddress string,
	constractAddress string,
	gas string,
	erc20Value string) error{

	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})

	res,err := req.Post(InsertReceivedERC20CoinInfoAPI,map[string]interface{} {
		"txHash":hash,
		"blockHash":blockHash,
		"blockNumber":blockNumber,
		"fromAddress":fromAddress,
		"toAddress":toAddress,
		"contractAddress":constractAddress,
		"erc20Value":erc20Value,
		"types":"BGFT",
		"gas":gas,
	})

	if nil != err {
		log.Println(err)
		return err
	}

	body,err := res.Body()
	if nil != err {
		log.Println(err)
		return err
	}

	resMap := make(map[string]interface{})
	err = json.Unmarshal(body,&resMap)
	if nil != err {
		log.Println(err)
		return err
	}


	stuCode,_  := resMap["status"].(float64)
	msg,_ := resMap["msg"].(string)

	if 0 == stuCode && 200 == res.StatusCode() {
		return nil
	}

	return errors.New(msg)

}

//根据交易Hash 确认交易
func (t EthTransactionHandler) NoticeTransAffirmed(transHash string) error {
	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})
	res,err := req.Post(NoticeTransSucceededAPI,map[string]interface{} {
		"txHash":transHash,
	})

	if nil != err {
		return err
	}

	body,err := res.Body()
	if nil != err {
		return err
	}

	resMap := make(map[string]interface{})
	err = json.Unmarshal(body,&resMap)
	if nil != err {
		return err
	}


	stuCode,_  := resMap["status"].(float64)
	msg,_ := resMap["msg"].(string)

	if 0 == stuCode && 200 == res.StatusCode() {
		return nil
	}


	return errors.New(msg)
}

// 根据交易Hash 重发交易
func (t EthTransactionHandler) NoticeTransFailed(transHash string) error{
	req := HttpRequest.NewRequest()

	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})
	res,err := req.Post(NoticeTransFailedAPI,map[string]interface{} {
		"txHash":transHash,
	})

	if nil != err {
		return err
	}

	body,err := res.Body()
	if nil != err {
		log.Println(err)
		return err
	}

	resMap := make(map[string]interface{})
	err = json.Unmarshal(body,&resMap)
	if nil != err {
		log.Println(err)
		return err
	}


	stuCode,_  := resMap["status"].(float64)
	msg,_ := resMap["msg"].(string)

	if 0 == stuCode && 200 == res.StatusCode() {
		return nil
	}

	return errors.New(msg)
}

func (t EthTransactionHandler) GetUnHandledTransInfo() []map[string]string {

	req := HttpRequest.NewRequest()
	res,err := req.Post(GetUnHandledTransInfoAPI,map[string]interface{} {

	})

	log.Println(res,err)
	return nil
}