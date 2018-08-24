package blocknode

import (
	"clmwallet-block-wacther/database/mysqltools"
	"errors"
)

type BlockNodeInfo struct {
	Number int64		`gorm:"primary_key"`
	Hash   string
	ParentHash string `gorm:"-"`
	TransHash string
}

func (BlockNodeInfo)  TableName() string	{
	return "blockNodeInfo"
}

func init() {
	db := mysqltools.GetInstance().GetMysqlDB()
	if !db.HasTable(&BlockNodeInfo{}) {
		db.CreateTable(&BlockNodeInfo{})
	}
}

func (info *BlockNodeInfo) Equal(info1 *BlockNodeInfo) bool {
	if nil == info1 {
		return false
	}
	return info.Number == info1.Number && info.Hash == info1.Hash
}

/// 持久化到数据库
func (info *BlockNodeInfo) Store()  error {
	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.Create(info).Error;err != nil {
		return err
	}
	return nil
}

/// 更新到数据库
func (info *BlockNodeInfo) Save() error  {
	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.Save(info).Error;err != nil {
		return err
	}
	return nil
}

/// 从数据库中删除信息
func (info *BlockNodeInfo) Delete() error  {
	if 0 >= info.Number {
		return errors.New("Primary key don't allow Empty.")
	}

	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.Delete(info).Error;err != nil {
		return err
	}
	return nil
}

func  Find(nodes *[]BlockNodeInfo) error  {
	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.Find(nodes).Error;err != nil {
		return err
	}
	return nil
}