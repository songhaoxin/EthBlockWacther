/*
@Time : 2018/9/1 下午5:23 
@Author : Mingjian Song
@File : configs
@Software: 深圳超链科技
*/

package configs

import (
	"io/ioutil"
	"encoding/json"
)

// geth结点地址
//const GethHost  = "http://localhost:7545"
//const GethHost  = "http://192.168.50.167:7545"
//const GethHost  = "https://ropsten.infura.io/v3/70905704824a4e66888eda02d084c49d"
const GethHost  = "http://47.75.115.210:8545"

// 管理平台服务地址
//const ServerHost  = "http://192.168.50.138:8781"//"http://47.75.115.210:8781"
const ServerDBConnectString  = "root:current@tcp(47.75.194.231)/multi_wallet?charset=utf8&parseTime=True&loc=Local"

// 保存同步区块数据库信息
//const BlockDataConnectString  = "root:root@tcp(127.0.0.1)/clwallet?charset=utf8&parseTime=True&loc=Local"
//const BlockDataConnectString  = "root:root@tcp(120.77.223.246)/clwallet?charset=utf8&parseTime=True&loc=Local"
const BlockDataConnectString  = "root:current@tcp(47.75.194.231)/multi_wallet?charset=utf8&parseTime=True&loc=Local"


// "47.106.136.96"
// 区块的确认高度
const AffiremBlockHeigh  = 6

//const MaxSizeAllowed  = 1000

// 定时时间片
const TimeDelayInSecand = 10

type Config struct {
	GethHost string `json:"GethHost"`
	ServerDBConnectString string `json:"ServerDBConnectString"`
	BlockDataConnectString string `json:"BlockDataConnectString"`
	AffiremBlockHeigh int `json:"AffiremBlockHeigh"`
	TimeDelayInSecand int `json:"TimeDelayInSecand"`
}

func Init() *Config  {
	return &Config{}
}

//var GobalConfig *Config = Init()

func (cfg *Config) Load(filename string) {
	data,err := ioutil.ReadFile(filename)
	if err != nil {
		cfg.GethHost = GethHost
		cfg.ServerDBConnectString = ServerDBConnectString
		cfg.BlockDataConnectString = BlockDataConnectString
		cfg.AffiremBlockHeigh = AffiremBlockHeigh
		cfg.TimeDelayInSecand = TimeDelayInSecand
		return
	}

	err = json.Unmarshal(data,cfg)
	if err != nil {
		cfg.GethHost = GethHost
		cfg.ServerDBConnectString = ServerDBConnectString
		cfg.BlockDataConnectString = BlockDataConnectString
		cfg.AffiremBlockHeigh = AffiremBlockHeigh
		cfg.TimeDelayInSecand = TimeDelayInSecand
	}
}
