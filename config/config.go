package config


// geth结点地址
const GethHost  = "http://localhost:7545"
//const GethHost  = "https://ropsten.infura.io/v3/70905704824a4e66888eda02d084c49d"

// 管理平台服务地址
const ServerHost  = "http://192.168.50.138:8781"//"http://47.75.115.210:8781"

// 保存同步区块数据库信息
//const BlockDataConnectString  = "root:root@tcp(127.0.0.1)/clwallet?charset=utf8&parseTime=True&loc=Local"
const BlockDataConnectString  = "root:root@tcp(120.77.223.246)/clwallet?charset=utf8&parseTime=True&loc=Local"


// 区块的确认高度
const AffiremBlockHeigh  = 10

const MaxSizeAllowed  = 1000

// 定时时间片
const TimeDelayInSecand = 3

