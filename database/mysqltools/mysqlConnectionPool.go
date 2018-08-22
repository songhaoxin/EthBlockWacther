package mysqltools

import (
	"log"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MysqlConnectionPool struct {
}

var instance *MysqlConnectionPool
var once sync.Once

var db *gorm.DB
var err_db error

func GetInstance() *MysqlConnectionPool {
	once.Do(func() {
		instance = &MysqlConnectionPool{}
		instance.InitDataPool()
	})

	return instance
}

func (m *MysqlConnectionPool) InitDataPool() (isSuccess bool) {

	db, err_db = gorm.Open("mysql", "root:root@tcp(120.77.223.246)/clwallet?charset=utf8&parseTime=True&loc=Local")
	if nil != err_db {
		log.Fatal(err_db)
		return false
	}

	return true

}

/// 对外获取数据库连接对象
func (m *MysqlConnectionPool) GetMysqlDB() (db_con *gorm.DB) {
	return db
}
