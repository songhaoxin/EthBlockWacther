package blocknode

import (
	"clmwallet-block-wacther/database/mysqltools"
	"errors"
)

type BlockNodeInfo struct {
	BlockNumber int64		`gorm:"primary_key"`
	BlockHash   string
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
	if 0 >= info.BlockNumber {
		return errors.New("Primary key don't allow Empty.")
	}

	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.Delete(info).Error;err != nil {
		return err
	}
	return nil
}

func FindAll()  {
	//nodeInfo := make([]BlockNodeInfo,2)
	//var nodes  []BlockNodeInfo
	//db := mysqltools.GetInstance().GetMysqlDB()
	//if err := db.Find(nodes).Error;err != nil {
	//
	//
	//}
	//
	//fmt.Println(nodeInfo)

}