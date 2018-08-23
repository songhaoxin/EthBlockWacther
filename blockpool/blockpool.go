package blockpool

import (
	"clmwallet-block-wacther/modles/blocknode"
	"sync"
	"clmwallet-block-wacther/config"
	"strings"
)

type BlockPool struct {
	MaxSizeAllowed    int64 //允许的区块最大数量
	AffiremBlockHeigh int64 //用以确认区块，离最新区块的标准高度

	startIdx int64            //区块池中的起始区块号
	endIdx   int64            //区块池中的最新区块号
	size     int64            //区块池中的区块数量
	pool     map[int64] *blocknode.BlockNodeInfo //管理最近区块的池子
	lock     *sync.RWMutex
}

/// 创建 "BlockPool"实例
func Init() *BlockPool {
	p := &BlockPool{
		startIdx:          -1,
		endIdx:            -1,
		size:              0,
		MaxSizeAllowed:    config.MaxSizeAllowed,
		AffiremBlockHeigh: config.AffiremBlockHeigh,
		pool:              make(map[int64] *blocknode.BlockNodeInfo),
		lock:              new(sync.RWMutex),
	}
	return p

}

func (b *BlockPool) IsEmpty() bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.size == 0
}

func (b *BlockPool) Size() int64 {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.size
}

func (b *BlockPool) EarliestNumber() int64 {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.startIdx
}

func (b *BlockPool) LatestNumber() int64 {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.endIdx
}

func (b *BlockPool) ContainElement(info *blocknode.BlockNodeInfo) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if nil == info || nil == b.pool{
		return false
	}

	if v, ok := b.pool[info.Number]; ok {
		if v.Hash == info.Hash {
			return true
		}
	}

	return false
}


/// 从数据库中全量加载所有记录
func (b *BlockPool) LoadBlocksFromDB()  {

}

/// 从数据库中加载一条记录
func (b *BlockPool) loadBlockFromDB(node *blocknode.BlockNodeInfo)  {
	if nil == node || nil == b.pool{
		return
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	if b.size + 1 == b.MaxSizeAllowed {
		b.removeElementAtStart()
	}

	k,v := node.Number,node

	if k < 0 {
		return
	}

	// 更新startIdx 与 endIdx
	if 0 == b.size {
		b.startIdx = k
		b.endIdx = k
	} else if k < b.startIdx {
		b.startIdx = k
	} else if k > b.endIdx {
		b.endIdx = k
	}

	b.pool[k] = v
	b.size++
}


/// 从区块链中接收一个区块信息，并找出孤立的区块（如果存在）
func (b *BlockPool) ReciveBlock(node *blocknode.BlockNodeInfo) *blocknode.BlockNodeInfo {
	if nil == node || nil == b.pool{
		return nil
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	if b.size+1 == b.MaxSizeAllowed {
		b.removeElementAtStart()
	}

	k, v := node.Number, node

	if k < 0  {
		return nil
	}

	// 更新startIdx 与 endIdx
	if 0 == b.size {
		b.startIdx = k
		b.endIdx = k
	} else if k < b.startIdx {
		b.startIdx = k
	} else if k > b.endIdx {
		b.endIdx = k
	}

	var n *blocknode.BlockNodeInfo = nil
	if ov,ok := b.pool[k]; ok { //存在旧值
		if ov.Hash != v.Hash { //Hash值不一样
			ov.Delete()
			n = ov
		} else { //Hash值一样,直接忽略
			return nil
		}
	} else {
		b.size++
	}

	// 更新或增加元素
	b.pool[k] = n
	node.Store()

	return n
}

/// 对区块进行校验，以处理已经被确认的交易
func (b *BlockPool) LookBocks4AffirmTrans() []string {

	b.lock.RLock()
	defer b.lock.RUnlock()

	if nil == b.pool {return nil}
	if b.size <= 0 {
		return  nil
	}

	affirmTransHashSlice := make([]string,0)

	for k, v := range b.pool {
		if b.endIdx-k+1 >= b.AffiremBlockHeigh {
			tHashs := strings.Split(v.TransHash, ";")
			for _,tHash := range tHashs {
				if "" != tHash {
					affirmTransHashSlice = append(affirmTransHashSlice,tHash)
				}
			}

		}
	}

	return  affirmTransHashSlice
}


func (b *BlockPool) removeElementAtStart() {
	b.lock.Lock()
	defer b.lock.Unlock()

	if nil == b.pool {return }

	if b.size == 0 {
		return
	}
	if b.startIdx < 0 {
		return
	}

	k := b.startIdx
	if v, ok := b.pool[k]; ok {
		delete(b.pool, k)
		v.Delete()
	}

	if 0 == b.size {
		b.startIdx = -1
		b.endIdx = -1
	} else {
		b.startIdx++
	}
	b.size--
}
