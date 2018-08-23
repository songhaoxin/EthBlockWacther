package config


// geth结点地址
const GethHost  = "http://localhost:7545"

// 保存同步区块数据库信息
//const BlockDataConnectString  = "root:root@tcp(127.0.0.1)/clwallet?charset=utf8&parseTime=True&loc=Local"
const BlockDataConnectString  = "root:root@tcp(120.77.223.246)/clwallet?charset=utf8&parseTime=True&loc=Local"
// 区块的确认高度
const AffiremBlockHeigh  = 5

const MaxSizeAllowed  = 1000

// 定时时间片
const TimeDelayInSecand = 10

