/*
@Time : 2018/9/1 上午11:14 
@Author : Mingjian Song
@File : ethhandler_test
@Software: 深圳超链科技
*/

package transactionhandler

import (
	"testing"
	"log"
	"fmt"
)

var hdl *EthTransactionHandler  = Init()

func TestExistAddress(t *testing.T)  {
	/*
	tcase := []struct {
		address string
		isExist bool
	}{
		// 本平台
		{"0x8fCC75bB4D90082D6EB3Aa9Da0d14DAbD538d34f",true},

		//不是本平台
		{"0x2DEF2400a4cC1aB21bfb6472B54A00643740E109",false},
	}

	for _,tt := range tcase{
		actual := hdl.ExistAddress(tt.address)
		if actual != tt.isExist {
			t.Errorf("测试用例：%s 的结果本来应该是 %v,但测试的结果却为:%v",tt.address,tt.isExist,actual)
		}
	}*/
	//for i:=1;i< 100000000000;i++{
		hdl.ExistAddress("0x8fCC75bB4D90082D6EB3Aa9Da0d14DAbD538d34f")
	//}

}

func BenchmarkEthTransactionHandler_ExistAddress(b *testing.B) {
	hdl.ExistAddress("0x2DEF2400a4cC1aB21bfb6472B54A00643740E109")
}

func TestExistTransByHash(t *testing.T)  {
	//ExistTransByHash
	tcase := []struct{
		hash string
		isExist bool
	}{
		// 本平台
		{hash:"0xa50a59ae8d4fcf672d0034d31cdaff48710b6bb5cac60c844693435b52f9dab2", isExist:true},
		// 不是本平台
		{"0x0091f3cdf24c393f20739acd499b76db891f1235",false},
	}

	for _,tt := range tcase {
		actual := hdl.ExistTransByHash(tt.hash)
		if actual != tt.isExist {
			t.Errorf("测试用例：%s 的结果本来应该是 %v,但其测试结果却为:%v",tt.hash,tt.isExist,actual)
		}
	}

}

func TestAddBlockNumberHash(t *testing.T)  {
	tcase := []struct{
		blockNumber int64
		blockHash string
		withTransHash string
		err error
	} {
		{
			9999,"hhhhhhhhhhsfsfsf","0x0ac50ee9ce8075be56aea8c893bdabaa96ec7f0c4345bdace93101ce75b4bd80",error(nil),
		},
		/*
		{
			"999","xxx","0x4dd0e3092ee7f397ac54975cfe7de374ae1eae124be131b45a9d9803ee03b4e5",error(nil),
		},*/
	}

	for _,tt := range tcase {
		actual := hdl.AddBlockNumberHash(tt.blockNumber,tt.blockHash,tt.withTransHash)
		if nil != actual {
			t.Errorf("测试用例：%s,%s,%s 的结果本来应该是 %v,但其测试结果却为:%v",tt.blockNumber,tt.blockHash,tt.withTransHash,tt.err,actual)
		}
	}
}

func TestInsertReceivedTransInfo(t *testing.T)  {

	tcase := []struct{
		hash string
		blockHash string
		blockNumber int64
		fromAddress string
		toAddress string
		gas string
		value string
		err error
	} {
		//成功的情况的测试用例
		{
			"33333","11",1111111111,"112","221","12","0x56bc75e2d6310000",nil,
		},
		/*
		//失败时候的测试用例
		{
			"","","","","","","",nil,
		},
*/
	}


	for _,tt := range tcase {
		actual := hdl.InsertReceivedTransInfo(tt.hash,tt.blockHash,tt.blockNumber,tt.fromAddress,tt.toAddress,tt.gas,tt.value)
		if nil != actual {
			//t.Errorf("测试用例：%s,%s,%s 的结果本来应该是 %v,但其测试结果却为:%v",tt.blockNumber,tt.blockHash,tt.withTransHash,tt.err,actual)
			t.Errorf("不通过")
		}
	}

}

func TestInsertReceivedERC20CoinInfo(t *testing.T)  {
	tcase := []struct{
		hash string
		blockHash string
		blockNumber int64
		fromAddress string
		toAddress string
		constractAddress string
		gas string
		erc20Value string
		err error
	} {
		//成功的情况的测试用例
		{
			"19961","12",4444444444444444,"23","32","0xd71c3ae0286286eac90dae97575d21c599ab0ffc11","232","0x56bc75e2d6310000",nil,
		},
		/*
		//失败时候的测试用例
		{
			"","","","","","","",nil,
		},
*/
	}


	for _,tt := range tcase {
		actual := hdl.InsertReceivedERC20CoinInfo(tt.hash,tt.blockHash,tt.blockNumber,tt.fromAddress,tt.toAddress,tt.constractAddress,tt.gas,tt.erc20Value)
		if nil != actual {
			//t.Errorf("测试用例：%s,%s,%s 的结果本来应该是 %v,但其测试结果却为:%v",tt.blockNumber,tt.blockHash,tt.withTransHash,tt.err,actual)
			t.Errorf("不通过")
		}
	}

}

func TestNoticeTransAffirmed(t *testing.T)  {
	tcase := []struct{
		txHash string
		err error
	} {
		// 应该成功的测试用例
		{"0x0ac50ee9ce8075be56aea8c893bdabaa96ec7f0c4345bdace93101ce75b4bd80",nil},
	}

	for _,tt := range tcase{
		actual := hdl.NoticeTransAffirmed(tt.txHash)
		if nil != actual {
			log.Println(actual)
			t.Errorf("测试用例：%s 的结果本来应该是正确的，但其测试结果却失败：%v",tt.txHash,actual)
		}
	}

}

func TestNoticeTransFailed(t *testing.T)  {
	tcase := []struct{
		txHash string
		err error
	} {
		// 应该成功的测试用例
		{"0x0ac50ee9ce8075be56aea8c893bdabaa96ec7f0c4345bdace93101ce75b4bd80",nil},
	}

	for _,tt := range tcase{
		actual := hdl.NoticeTransFailed(tt.txHash)
		if nil != actual {
			log.Println(actual)
			t.Errorf("测试用例：%s 的结果本来应该是正确的，但其测试结果却失败：%v",tt.txHash,actual)
		}
	}

}

func TestEthTransactionHandler_UpdateTransState(t *testing.T) {
	err := hdl.UpdateTransState("0x0ac50ee9ce8075be56aea8c893bdabaa96ec7f0c4345bdace93101ce75b4bd80",1)
	if nil != err {
		log.Println(err)
	}
}


func TestEthTransactionHandler_GetTransMainInfoes(t *testing.T) {
	f,to,a,con,e := hdl.GetTransMainInfo("0x0ac50ee9ce8075be56aea8c893bdabaa96ec7f0c4345bdace93101ce75b4bd80")

	fmt.Println(f,to,a,con,e)
}

func TestEthTransactionHandler_sendMessage(t *testing.T) {
	hdl.sendMessage("0xebee9bfbdf9427a4be43671056cae7cdc601e53b307e0d524d0704e4a43a4cd8",1)
}

func TestEthTransactionHandler_SetLowestIdxFromGether(t *testing.T) {
	if nil != hdl.SetLowestIdxFromGether(120) {
		t.Errorf("错误")
	}
}

func TestEthTransactionHandler_GetLowestIdxFromServer(t *testing.T) {
	n := hdl.GetLowestIdxFromServer()
	fmt.Println(n)
}







