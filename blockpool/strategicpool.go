/*
@Time : 2018/8/30 下午8:58
@Author : Mingjian Song
@File : strategicpool
@Software: 深圳超链科技
*/

package blockpool

import (
	"clmwallet-block-wacther/configs"
	"clmwallet-block-wacther/modles/blocknode"
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/metrics"
	"strings"
)

type StrategicPool struct {
	affiremHeigh int64 // 用以确认区块的高度
	//earliestIdx		int64
	latestIdx int64 // 拉取的最新区块号
	size      int
	pool      map[int64]*blocknode.BlockNodeInfo
	lock      *sync.RWMutex
}

func Init() *StrategicPool {
	p := &StrategicPool{
		affiremHeigh: configs.AffiremBlockHeigh,
		latestIdx:    -1,
		size:         0,
		//earliestIdx:-1,
		pool: make(map[int64]*blocknode.BlockNodeInfo),
		lock: new(sync.RWMutex),
	}
	p.LoadBlocksFromDB()
	return p
}

/// 从数据库中全量加载所有记录
func (s *StrategicPool) LoadBlocksFromDB() {
	var nodes []blocknode.BlockNodeInfo
	if err := blocknode.Find(&nodes); err != nil {
		return
	}
	for _, n := range nodes {
		var node blocknode.BlockNodeInfo = n
		s.InsertElement(&node)
	}
}

func (s *StrategicPool) InsertElement(node *blocknode.BlockNodeInfo) {
	if s.ContainElement(node) {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	k, v := node.Number, node

	if k < 0 {
		log.Println("非法区块号（负数），忽略该区块")
		return
	}

	s.pool[k] = v
	s.size++
}

func (s *StrategicPool) Save2Db() {
	s.lock.Lock()
	defer s.lock.Unlock()

	log.Println("保存数据到数据库...")
	for _, v := range s.pool {
		v.Save()
	}
	log.Println("保存数据成功！")
}

/// 从区块链中接收一个区块信息，并找出孤立的区块（如果存在）
func (s *StrategicPool) ReciveBlockFromChain(node *blocknode.BlockNodeInfo) *blocknode.BlockNodeInfo {

	fmt.Println(node.Number)
	if s.ContainElement(node) {
		log.Println("本地已经存在该区块，跳过！")
		return nil
	}

	// 测试时临时注释
	if "" == node.TransHash {
		log.Println("不包括本平台帐户的交易信息，跳过！")
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	k, v := node.Number, node

	if k < 0 {
		log.Println("区块号出现负数，忽略该区块!")
		return nil
	}

	var n *blocknode.BlockNodeInfo = nil
	if ov, ok := s.pool[k]; ok { //存在旧值
		if ov.Hash != v.Hash { //Hash值不一样
			ov.Delete()
			n = ov
		}
	} else {
		s.size++
	}

	// 更新或增加元素
	s.pool[k] = node

	log.Printf("更新了池中数据")
	node.Save()
	log.Println("保存到了数据库中了")

	return n
}

// 对区块进行校验，以处理已经被确认的交易
func (s *StrategicPool) LookSuccessedTransHashs() []string {

	s.lock.Lock()
	defer s.lock.Unlock()

	log.Println("Begin scan block for suceess 1...")

	if 0 == s.size || s.latestIdx < 0 {
		return nil
	}

	log.Println("Begin scan block for suceess 2...")

	affirmTransHashSlice := make([]string, 0)

	var keys metrics.Int64Slice
	for k := range s.pool {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	var k int64
	for _, k = range keys {

		if s.latestIdx-k < int64(s.affiremHeigh) {
			continue
		}

		v, ok := s.pool[k]
		if ok {
			tHashs := strings.Split(v.TransHash, ";")
			for _, tHash := range tHashs {
				if "" != tHash {
					affirmTransHashSlice = append(affirmTransHashSlice, tHash)
					log.Println("交易成功的区块号：", v.Number)
					log.Println("交易成功的区块HASH:", tHash)
				}
			}

			//从池中删除这个区块
			delete(s.pool, v.Number)
			s.size--
			v.Delete()
		}

	}

	return affirmTransHashSlice
}

func (s *StrategicPool) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.size
}


func (s *StrategicPool) ContainElement(info *blocknode.BlockNodeInfo) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if nil == info || nil == s.pool {
		return false
	}

	if v, ok := s.pool[info.Number]; ok {
		if v.Hash == info.Hash {
			return true
		}
	}

	return false
}

func (s *StrategicPool) Descrip() {
	fmt.Println("Size:", s.Size())
	fmt.Println("earliestIdx:", s.GetEarliestIdx())
	fmt.Println("affiremHeigh:", s.GetAffiremHeigh())
	fmt.Println("latestIdx:", s.GetLatestIdx())

	for k, v := range s.pool {
		fmt.Println(k, v)
	}
}

func (s *StrategicPool) GetEarliestIdx() int64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if 0 == s.size {
		return -1
	}

	var keys metrics.Int64Slice
	for k := range s.pool {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	return keys[0]
}

func (s *StrategicPool) GetAffiremHeigh() int64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.affiremHeigh
}

func (s *StrategicPool) GetLatestIdx() int64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.latestIdx
}

func (s *StrategicPool) SetLatestIdx(idx int64) {
	if 0 > idx {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	s.latestIdx = idx
}
