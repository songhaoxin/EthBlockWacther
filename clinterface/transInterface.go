package clinterface

type TransInterface interface {
	//根据传入帐号地址，判断是否属于超链钱包平台的
	ExistAddress(address string) bool

	//根据传入的交易Hash，判断是否已经存在
	ExistTransByHash(transHash string) bool

	//根据交易Hash 增加blockNumber/blockHash到交易数据库表
	AddBlockNumberHash(blockNumber string,blockHash string,withTransHash string)

	//根据解析的交易的信息（别人向我们的帐户转帐这一情况），添加一条记录到交易表
	InsertTransInfo(
		fromAddress string,
		toAddress string,
		value string,
		gas string,
		txHash string)

	//确认交易 （区块号）

	// 重发交易 (区块号）
}
