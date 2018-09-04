package transactionhandler

import (
	"github.com/kirinlabs/HttpRequest"
	"log"
	"clmwallet-block-wacther/configs"
	"errors"

	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"clmwallet-block-wacther/helper"
	"encoding/json"
)



type EthTransactionHandler struct {
	db *sql.DB
}


func Init() *EthTransactionHandler {
	return &EthTransactionHandler{}
}

func (e EthTransactionHandler) CloseDB()  {
	if nil != e.db {
		e.db.Close()
	}
}

///====================== TransInterface协议方法========================
//////////////////////////////////////////////////////////////////////////////////
func (e EthTransactionHandler) getMySqlDb() *sql.DB  {
	var err error = nil
	if nil == e.db {
		e.db,err = sql.Open("mysql",configs.ServerDBConnectString)
		if nil != err {
			log.Println(err)
			return nil
		}
	}
	return e.db
}


// 判断指定的 地址 是否属于平台内帐户
func (t EthTransactionHandler) ExistAddress(address string) bool {

	db := t.getMySqlDb()
	if nil == db {
		return false
	}

	var address2 string
	err := db.QueryRow("SELECT address FROM coin_address WHERE address=?",address).Scan(&address2)
	if nil != err {
		log.Println(err)
		return false
	}

	return true

}

func (t EthTransactionHandler) ExistTransByHash(transHash string) bool {
	db := t.getMySqlDb()
	if nil == db {
		return false
	}

	var transHash2 string
	err := db.QueryRow("SELECT tx_hash FROM send_out_detail WHERE tx_hash=?",transHash).Scan(&transHash2)
	if nil != err {
		log.Println(err)
		return false
	}


	return true
	
}

//根据交易Hash 增加blockNumber/blockHash到交易数据库表
func (t EthTransactionHandler) AddBlockNumberHash(blockNumber string,blockHash string,withTransHash string) error{

	db := t.getMySqlDb()
	if nil == db {
		return errors.New("Can not conect to DB Server")
	}

	stmt,_ := db.Prepare("UPDATE send_out_detail SET block_number=?,block_hash=? WHERE tx_hash=?")
	defer stmt.Close()

	ret,err := stmt.Exec(blockNumber,blockHash,withTransHash)
	if nil != err {
		log.Println("AddBlockNumberHash: insert data error:%v",err)
		return err
	}

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println("RowsAffected:", RowsAffected)
	}

	return nil

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

	db := t.getMySqlDb()
	if nil == db {
		return errors.New("Can not conect to DB Server")
	}

	stmt,_ := db.Prepare("INSERT INTO send_out_detail (tx_hash,block_hash,block_number,from_address,to_address,gas,to_amount) VALUES (?,?,?,?,?,?,?)")
	defer stmt.Close()


	ret,err := stmt.Exec(hash,blockHash,blockNumber,fromAddress,toAddress,gas,value)
	if nil != err {
		log.Println("InsertReceivedTransInfo: insert data error:%v",err)
		return err
	}

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println("RowsAffected:", RowsAffected)
	}

	return nil

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

	db := t.getMySqlDb()
	if nil == db {
		return errors.New("Can not conect to DB Server")
	}

	stmt,_ := db.Prepare("INSERT INTO send_out_detail (tx_hash,block_hash,block_number,from_address,to_address,contract_address,gas,to_amount) VALUES (?,?,?,?,?,?,?,?)")
	defer stmt.Close()


	ret,err := stmt.Exec(hash,blockHash,blockNumber,fromAddress,toAddress,constractAddress,gas,erc20Value)
	if nil != err {
		log.Println("InsertReceivedERC20CoinInfo: insert data error:%v",err)
		return err
	}

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println("RowsAffected:", RowsAffected)
	}

	return nil


}



//根据交易Hash 确认交易
func (e EthTransactionHandler) NoticeTransAffirmed(transHash string) error {
	log.Println("发：交易成功 短信------->>>>成功 ---HASH:%s",transHash)
	return e.sendMessage(transHash,1)
}

// 根据交易Hash 重发交易
func (e EthTransactionHandler) NoticeTransFailed(transHash string) error{
	log.Println("发：交易失败 短信------->>>>失败 ---HASH:%s",transHash)
	return e.sendMessage(transHash,2)
}



///私有方法
///////////////////////////////////////////////////////////////////////////////////////////////////////
func (e EthTransactionHandler)sendMessage(transHash string,stateTag int) error {
	fromAddress,toAddress,amount,contract,err := e.GetTransMainInfo(transHash)
	if nil != err {return err}

	if "" == fromAddress {
		log.Println("Error:fromAddress is empty!")
		return errors.New("Error:fromAddress is empty!")
	}

	if "" == toAddress {
		log.Println("Error:toAddress is empty!")
		return errors.New("Error:toAddress is empty!")
	}

	if "" == amount {
		log.Println("Error:amount is empty!")
		return errors.New("Error:amount is empty!")
	}


	fromPhoneNumber := e.PhoneNumber(fromAddress)
	toPhoneNumber := e.PhoneNumber(toAddress)
	smbl,dec,err := e.GetCoinInfo(contract)
	if nil != err {
		log.Println(err)
		return err
	}

	amountDec,err := helper.Hex2Decimal(amount,dec,6)
	if nil != err {
		log.Println(err)
		return err
	}
	amount = amountDec.String() + smbl

	if 1 == stateTag { /// 更新交易的确认状态为确认
		e.UpdateTransState(transHash)
	} else { /// 删除未成功的交易记录
		// 暂未实现
	}

	///////////////////////////////////////////////////////////////////////
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})

	toAddress2 := helper.Substr(toAddress,0,5) + "..." + helper.Substr2(toAddress,37,41)

	res,err := req.Post("http://api.bgft.ltd/api/yzm/send_code",map[string]interface{} {
		"from":fromPhoneNumber,//13926514670,
		"to":toPhoneNumber,
		"to_address":toAddress2,
		"money":amount,
		"type":stateTag,
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
	log.Println(resMap)

	return nil
}



func (e EthTransactionHandler)PhoneNumber(address string) string  {
	db := e.getMySqlDb()
	if nil == db {
		return ""
	}

	var phoneNumber string
	err := db.QueryRow("SELECT phone_num FROM wallet WHERE id IN (SELECT wallet_id_id FROM coin_address WHERE address=?)",address).Scan(&phoneNumber)
	if nil != err {
		log.Println(err)
		return ""
	}
	return phoneNumber
}

func (t EthTransactionHandler)GetCoinInfo(contract string) (string,int,error) {
	db := t.getMySqlDb()
	if nil == db {
		return "",0,errors.New("Can not conect to DB Server")
	}

	var symbol string
	var decimal int
	err := db.QueryRow("SELECT `mark`,`decimal` FROM tokens WHERE contract=?",contract).Scan(&symbol,&decimal)
	if nil != err {
		return "",0,err
	}

	return symbol,decimal,nil
}

func (t EthTransactionHandler)UpdateTransState(hash string) error  {
	db := t.getMySqlDb()
	if nil == db {
		return errors.New("Can not conect to DB Server")
	}

	stmt,_ := db.Prepare("UPDATE send_out_detail SET is_confirm=? WHERE tx_hash=?")
	defer stmt.Close()

	ret,err := stmt.Exec(1,hash)
	if nil != err {
		log.Println("UpdateTransState: set  is_confirm error:%v",err)
		return err
	}

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println("LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println("RowsAffected:", RowsAffected)
	}

	return nil

}




func (t EthTransactionHandler)GetTransMainInfo(hash string) ( string, string, string,string, error)  {
	db := t.getMySqlDb()
	if nil == db {
		return "","","","",errors.New("Can not conect to DB Server")
	}

	var fromAddress string
	var toAddress string
	var amount string
	var contract string


	err:= db.QueryRow("SELECT from_address,to_address,to_amount,contract_address FROM send_out_detail WHERE tx_hash=?",hash).Scan(&fromAddress,&toAddress,&amount,&contract)
	if nil != err {
		//log.Println(err)
		return "","","","",err
	}

	return fromAddress,toAddress,amount,contract,nil

}

