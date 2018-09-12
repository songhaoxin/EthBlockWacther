package clinterface

type TransInterface interface {
	//根据传入帐号地址，判断是否属于超链钱包平台的
	ExistAddress(address string) bool

	//根据传入的交易Hash，判断是否已经存在
	ExistTransByHash(transHash string) bool

	//根据交易Hash 增加blockNumber/blockHash到交易数据库表
	AddBlockNumberHash(blockNumber int64,blockHash string,withTransHash string) error

	//根据解析的交易的信息（别人向我们的帐户转帐这一情况），添加一条记录到交易表
	InsertReceivedTransInfo(
		hash string,
		blockHash string,
		blockNumber int64,
		fromAddress string,
		toAddress string,
		gas string,
		value string) error

	//根据解析的ERC20代币交易信息(别人向我们的帐户转帐这一情况），添加一条记录到交易表
	InsertReceivedERC20CoinInfo(
		hash string,
		blockHash string,
		blockNumber int64,
		fromAddress string,
		toAddress string,
		constractAddress string,
		gas string,
		erc20Value string) error

	//根据交易Hash 确认交易
	NoticeTransAffirmed(transHash string) error

	// 根据交易Hash 重发交易
	NoticeTransFailed(transHash string) error


	// 返回平台中发出去(或者别人转给我们的）， --------> 如果收到的交易和发出去的交易分别设计表的话，需要操作两张表
	// 但是还没有被确认的交易的信息
	// 用于当监视服务挂掉后重启时，追溯需要确认的交易信息
	// 返回的格式
	// [
	//  { "blockNumber":"xxxxxxxx",
	//   "blockHash":"xxxxxx"
	//   "transactions":"transHash1;transHash2;transHash3"                ------->一个区块中有多条交易时用";"分隔
	//  }，
	//  {....},
	//  {....}
	// ]
	//
	//包括所以太坊币交易和ERC20代币交易
	//GetUnHandledTransInfo() []map[string]string



	GetLatestNumberShould2Fecth() int64


}
