/*
@Time : 2018/9/1 下午5:23 
@Author : Mingjian Song
@File : configs
@Software: 深圳超链科技
*/

package configs



// geth结点地址
//const GethHost  = "http://localhost:7545"
//const GethHost  = "http://192.168.50.167:7545"
const GethHost  = "https://ropsten.infura.io/v3/70905704824a4e66888eda02d084c49d"

// 管理平台服务地址
//const ServerHost  = "http://192.168.50.138:8781"//"http://47.75.115.210:8781"
const ServerDBConnectString  = "root:current@tcp(47.75.194.231)/multi_wallet?charset=utf8&parseTime=True&loc=Local"

// 保存同步区块数据库信息
//const BlockDataConnectString  = "root:root@tcp(127.0.0.1)/clwallet?charset=utf8&parseTime=True&loc=Local"
const BlockDataConnectString  = "root:root@tcp(120.77.223.246)/clwallet?charset=utf8&parseTime=True&loc=Local"
//const BlockDataConnectString  = "root:current@tcp(47.106.136.96)/multi_wallet?charset=utf8&parseTime=True&loc=Local"


// "47.106.136.96"
// 区块的确认高度
const AffiremBlockHeigh  = 3

//const MaxSizeAllowed  = 1000

// 定时时间片
const TimeDelayInSecand = 10

type Config struct {
	GethHost string `json:"geth_host"`
	ServerHost string `json:"server_host"`
	BlockDataConnectString string `json:"block_data_connect_string"`
	AffiremBlockHeigh int `json:"affirem_block_heigh"`
	TimeDelayInSecand int `json:"time_delay_in_secand"`
}
