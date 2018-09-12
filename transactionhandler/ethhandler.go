package transactionhandler

import (
	"log"
	"clmwallet-block-wacther/configs"
	"errors"

	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"strings"
	"github.com/kirinlabs/HttpRequest"
	"clmwallet-block-wacther/helper"
	"encoding/json"
)



type EthTransactionHandler struct {
	db *sql.DB
}


func Init() *EthTransactionHandler {
	return &EthTransactionHandler{}
}

func (e *EthTransactionHandler) CloseDB()  {
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

	/*
	db := t.getMySqlDb()
	if nil == db {
		return false
	}

	var address2 string

	err := db.QueryRow("SELECT address FROM coin_address WHERE address=?",address).Scan(&address2)
	if nil != err {
		//log.Println(err,"ExistAddress")
		return false
	}
	return true
*/

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return false
	}
	defer db.Close()

	var address2 string

	err = db.QueryRow("SELECT address FROM coin_address WHERE address=?",address).Scan(&address2)
	if nil != err {
		//log.Println(err,"ExistAddress")
		return false
	}
	return true



}

func (t EthTransactionHandler) ExistTransByHash(transHash string) bool {

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return false
	}
	defer db.Close()

	var transHash2 string
	err = db.QueryRow("SELECT tx_hash FROM send_out_detail WHERE tx_hash=?",transHash).Scan(&transHash2)
	if nil != err {
		log.Println(err,"ExistTransByHash")
		return false
	}


	return true
	
}

//根据交易Hash 增加blockNumber/blockHash到交易数据库表
func (t EthTransactionHandler) AddBlockNumberHash(blockNumber int64,blockHash string,withTransHash string) error{

	//db := t.getMySqlDb()
	//if nil == db {
	//	return errors.New("Can not conect to DB Server")
	//}

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return errors.New("Can not conect to DB Server")
	}
	defer db.Close()

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
	blockNumber int64,
	fromAddress string,
	toAddress string,
	gas string,
	value string) error {

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return errors.New("Can not conect to DB Server")
	}
	defer db.Close()

	//把value转换成十进制数
	valueDec,err := helper.Hex2Decimal(value,18,8)
	if nil != err {
		log.Println(err)
		return err
	}
	value = valueDec.String()


	stmt,_ := db.Prepare("INSERT INTO send_out_detail (tx_hash,block_hash,block_number,from_address,to_address,gas,to_amount,tx_type) VALUES (?,?,?,?,?,?,?,?)")
	defer stmt.Close()


	ret,err := stmt.Exec(hash,blockHash,blockNumber,fromAddress,toAddress,gas,value,1)
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
func (e EthTransactionHandler)  InsertReceivedERC20CoinInfo(
	hash string,
	blockHash string,
	blockNumber int64,
	fromAddress string,
	toAddress string,
	constractAddress string,
	gas string,
	erc20Value string) error{

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return errors.New("Can not conect to DB Server")
	}
	defer db.Close()

	_,dec,err := e.GetCoinInfo(constractAddress)
	if nil != err {
		log.Println(err)
		return err
	}
	log.Println("decimal:",dec)

	//把value转换成十进制数
	erc20Dec,err := helper.Hex2Decimal(erc20Value,dec,8)
	if nil != err {
		log.Println(err)
		return err
	}
	erc20Value = erc20Dec.String()

	stmt,_ := db.Prepare("INSERT INTO send_out_detail (tx_hash,block_hash,block_number,from_address,to_address,contract_address,gas,to_amount,tx_type) VALUES (?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()


	ret,err := stmt.Exec(hash,blockHash,blockNumber,fromAddress,toAddress,constractAddress,gas,erc20Value,1)
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

	if !e.ExistTransByHash(transHash) {
		log.Println("HASH错误，不是本平台的交易HASH")
	}

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
	fmt.Println("fromphone:",fromPhoneNumber)
	toPhoneNumber := e.PhoneNumber(toAddress)
	fmt.Println("tophone:",toPhoneNumber)

	var smbl string
	var dec int


	contract = strings.Replace(contract, " ", "", -1)
	if "" != contract { //是代币交易
		smbl,dec,err = e.GetCoinInfo(contract)
		if nil != err {
			log.Println(err)
			return err
		}
	} else {
		smbl = "ETH"
		dec = 18
	}
	fmt.Println(smbl,dec)


	amount = amount + smbl
	fmt.Println(amount)

	e.UpdateTransState(transHash,stateTag)


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

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return ""
	}
	defer db.Close()

	var phoneNumber string
	err = db.QueryRow("SELECT phone_num FROM wallet WHERE id IN (SELECT wallet_id_id FROM coin_address WHERE address=?)",address).Scan(&phoneNumber)
	if nil != err {
		log.Println(err)
		return ""
	}
	return phoneNumber
}

func (t EthTransactionHandler)GetCoinInfo(contract string) (string,int,error) {

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return "",0,errors.New("Can not conect to DB Server")
	}
	defer db.Close()

	var symbol string
	var decimal int
	err = db.QueryRow("SELECT `mark`,`decimal` FROM tokens WHERE contract=?",contract).Scan(&symbol,&decimal)
	if nil != err {
		return "",0,err
	}

	return symbol,decimal,nil
}

func (t EthTransactionHandler)UpdateTransState(hash string,state int) error  {
	if 1 != state && 2 != state {
		return errors.New("Paramas `state` allowed range of  1 or 2")
	}

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return errors.New("Can not conect to DB Server")
	}
	defer db.Close()

	stmt,_ := db.Prepare("UPDATE send_out_detail SET is_confirm=? WHERE tx_hash=?")
	defer stmt.Close()

	ret,err := stmt.Exec(state,hash)
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




func (e EthTransactionHandler)GetTransMainInfo(hash string) ( string, string, string,string, error)  {

	db,err := sql.Open("mysql",configs.ServerDBConnectString)
	if nil != err {
		return "","","","",errors.New("Can not conect to DB Server")
	}
	defer db.Close()

	var fromAddress string
	var toAddress string
	var amount string
	var contract string


	err= db.QueryRow("SELECT from_address,to_address,to_amount,contract_address FROM send_out_detail WHERE tx_hash=?",hash).Scan(&fromAddress,&toAddress,&amount,&contract)
	if nil != err {
		return "","","","",err
	}

	return fromAddress,toAddress,amount,contract,nil

}



/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t EthTransactionHandler)GetLatestNumberShould2Fecth() int64  {
	/*
	number1 := t.GetLowestIdxFromGether()
	number2 := t.GetLowestIdxFromServer()
	if number1 < number2 {
		if number1 >= 0 {
			return number1
		} else if number2 >= 0 {
			return number2
		}
	} else {
		if number2 <= 0 {
			return number1
		} else {
			return number2
		}
	}
	*/
	return -1
}

//从服务端获取的 需要确认的最早的区块号
func (t EthTransactionHandler) GetLowestIdxFromServer() int64 {
	/*
	db, err := sql.Open("mysql", configs.ServerDBConnectString)
	if nil != err {
		return -1, err
	}
	defer db.Close()
	*/

	db := t.getMySqlDb()
	if nil == db {
		return -1
	}

	var lowestNumber int64
	err := db.QueryRow(`SELECT  block_number FROM send_out_detail  WHERE is_confirm = 0 ORDER BY block_number ASC LIMIT 0,1`).Scan(&lowestNumber)
	if nil != err {
		return -1
	}
	if 0 == lowestNumber {
		lowestNumber = -1
	}
	return lowestNumber
}


/*
func (t *EthTransactionHandler) SetLowestIdxFromGether(number int64) error {

	idx := t.GetLowestIdxFromGether()
	db := t.getMySqlDb()
	if idx != -1 {
		stmt,_ := db.Prepare("UPDATE block_fecthed_info SET latestnumber=?")
		defer stmt.Close()

		_,err := stmt.Exec(number)
		if nil != err {
			log.Println("SetLowestIdxFromGether(ADD): ADD data error:%v",err)
			return err
		}
	} else {
		stmt,_ := db.Prepare("INSERT INTO block_fecthed_info (latestnumber) VALUES (?)")
		defer stmt.Close()

		_,err := stmt.Exec(number)
		if nil != err {
			log.Println("SetLowestIdxFromGether(UPDATE): UPDATE data error:%v",err)
			return err
		}
	}

	return nil
}
*/

//从服务端获取的 需要确认的最早的区块号
func (t EthTransactionHandler) GetLowestIdxFromGether() int64 {
	db := t.getMySqlDb()
	if nil == db {
		return -1
	}

	var lowestNumber int64
	err := db.QueryRow(`SELECT latestnumber FROM block_fecthed_info `).Scan(&lowestNumber)
	if nil != err {
		return -1
	}

	return lowestNumber
}


