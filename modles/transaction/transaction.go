package transaction

import (
	"github.com/jinzhu/gorm"
	"clmwallet/database/mysqltools"
)


type TransactionInfo struct {
	gorm.Model
	FromAddress string `gorm:size:60`
	ToAddress string	`gorm:size:60`
	BlockNumber string  `gorm:size:60`
	TxHash string		`gorm:size:80`
	Gas string			`gorm:size:30`
	Value string		`gorm:size:16`
	Input string		`gorm:size:16`
	Remark string
	Affirem bool		`gorm:default:false`
}

func (TransactionInfo)  TableName() string	{
	return "transactionInfo"
}

func init() {
	db := mysqltools.GetInstance().GetMysqlDB()
	if !db.HasTable(&TransactionInfo{}) {
		/*
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&TransactionInfo{}).Error;err != nil {
			panic(err)
		}*/
		db.CreateTable(&TransactionInfo{})
	}
}




/// 持久化到数据库
func (info *TransactionInfo) Store()  error {
	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.Create(info).Error;err != nil {
		return err
	}
	return nil
}

/// 更新到数据库
func (info *TransactionInfo) Save() error  {
	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.Save(info).Error;err != nil {
		return err
	}
	return nil
}


func GetAll()  {

}
