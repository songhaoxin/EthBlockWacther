package clinterface

type FecthBlockInterface interface {

	//根据区块号获取区块信息
	FecthBlockByNumber(blockNumber string) (info map[string]interface{},err error)

	//根据区块Hash获取区块信息
	FecthBlockByHash(blockHash string) (info map[string]interface{},err error)

}
