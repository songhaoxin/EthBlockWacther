/*
@Time : 2018/9/1 上午11:38 
@Author : Mingjian Song
@File : strategicpool_test
@Software: 深圳超链科技
*/

package blockpool

import (
	"testing"
	"sort"
	"strings"
	"log"
	"github.com/ethereum/go-ethereum/metrics"
	"fmt"
)

func TestStrategicPool_GetEarliestIdx(t *testing.T) {
	s := Init()
	fmt.Println(s.GetEarliestIdx())
	fmt.Println(s.Size())
	//fmt.Println(s.GetAffiremHeigh())

	fmt.Println("---------------------")

	//v,_ := s.pool[1]
	//v.Delete()

	delete(s.pool,1)

	delete(s.pool,59)
	delete(s.pool,60)
	fmt.Println(s.GetEarliestIdx())
	fmt.Println(s.Size())


}
func TestResetEarliestIdx(t *testing.T)  {

	s := Init()

	s.lock.Lock()
	defer s.lock.Unlock()

	if 0 == s.size {return }

	affirmTransHashSlice := make([]string,0)
	var keys metrics.Int64Slice
	for k := range s.pool {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	var k int64


	s.latestIdx = 4
	for _, k = range keys {
		fmt.Println(s.latestIdx)
		if s.latestIdx - k < int64(s.affiremHeigh) {
			continue
		}

		v,ok := s.pool[k]
		if ok {
			tHashs := strings.Split(v.TransHash, ";")
			for _,tHash := range tHashs {
				if "" != tHash {
					affirmTransHashSlice = append(affirmTransHashSlice,tHash)
					log.Println("交易成功的区块号：",v.Number)
					log.Println("交易成功的区块HASH:",tHash)
				}
			}

			//从池中删除这个区块
			delete(s.pool,v.Number)
			//v.Delete()

		}
	}

}
