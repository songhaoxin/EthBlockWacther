package blocknode

import (
	"clmwallet-block-wacther/database/mysqltools"
	"errors"
	"sync"
	"log"
)

type BlockNodeInfo struct {
	Number int64		`gorm:"primary_key"`
	Hash   string
	ParentHash string `gorm:"-"`
	TransHash string
	rwLock *sync.RWMutex `gorm:"-"`
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

func (info *BlockNodeInfo)getRWLock() *sync.RWMutex {
	if nil == info.rwLock {
		info.rwLock = new(sync.RWMutex)
	}
	return info.rwLock

}
func (info *BlockNodeInfo) Equal(info1 *BlockNodeInfo) bool {
	if nil == info1 {
		return false
	}
	return info.Number == info1.Number && info.Hash == info1.Hash
}

func (info *BlockNodeInfo)Exist() bool {
	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.First(info).Error;nil != err {
		log.Println(err)
		return false
	}

	return true
}

/// 保存或更新数据
func (info *BlockNodeInfo)Save() error  {
	if info.Number < 0 {
		return errors.New("区块号不能为负数")
	}
	if exist(info.Number) {
		info.update()
	} else {
		info.store()
	}

	return nil
}

func exist(number int64) bool  {
	var node = &BlockNodeInfo{}
	db := mysqltools.GetInstance().GetMysqlDB()
	db.First(node, number) // 查询number为number的node

	return "" != node.Hash
}

/// 持久化到数据库
func (info *BlockNodeInfo) store()  error {
	db := mysqltools.GetInstance().GetMysqlDB()
	if err := db.Create(info).Error;err != nil {
		return err
	}
	return nil
}

/// 更新到数据库
func (info *BlockNodeInfo) update() error  {
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