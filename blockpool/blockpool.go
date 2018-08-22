package blockpool

import (
	"clmwallet-block-wacther/modles/blocknode"
	"sync"
)

type BlockPool struct {
	MaxSizeAllowed    int64 //允许的区块最大数量
	AffiremBlockHeigh int64 //用以确认区块，离最新区块的标准高度

	startIdx int64            //区块池中的起始区块号
	endIdx   int64            //区块池中的最新区块号
	size     int64            //区块池中的区块数量
	pool     map[int64]string //管理最近区块的池子
	lock     *sync.RWMutex
}

/// 创建 "BlockPool"实例
func InitBlockPool() *BlockPool {
	p := &BlockPool{
		startIdx:          -1,
		endIdx:            -1,
		size:              0,
		MaxSizeAllowed:    100,
		AffiremBlockHeigh: 6,
		pool:              make(map[int64]string),
		lock:              new(sync.RWMutex),
	}
	return p

}

func (b *BlockPool) IsEmpty() bool {
	return b.size == 0
}

func (b *BlockPool) LatestNumber() int64 {
	return b.endIdx
}



func (b *BlockPool) ReciveBlock(node *blocknode.BlockNodeInfo) *blocknode.BlockNodeInfo {
	if nil == node {
		return nil
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	if b.size+1 == b.MaxSizeAllowed {
		b.removeElementAtStart()
	}

	k, v := node.BlockNumber, node.BlockHash

	if k < 0 || 0 == len(v) {
		return nil
	}

	// 更新startIdx 与 endIdx
	if k < b.startIdx {
		b.startIdx = k
	} else if k > b.endIdx {
		b.endIdx = k
	}

	var n *blocknode.BlockNodeInfo = nil
	if oldHash, ok := b.pool[k]; ok { //存在旧值
		if oldHash != v {
			n = &blocknode.BlockNodeInfo{k, oldHash,"",""}
		}
	} else {
		b.size++
	}

	// 更新或增加元素
	b.pool[k] = v
	return n
}

/// 对区块进行校验，以处理已经被确认的交易
func (b *BlockPool) LookBocks4AffirmTrans() (bool, *blocknode.BlockNodeInfo) {

	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.size <= 0 {
		return false, nil
	}

	for k, v := range b.pool {
		if b.endIdx-k+1 >= b.AffiremBlockHeigh {
			return true, &blocknode.BlockNodeInfo{k, v,"",""}
		}
	}
	return false, nil
}

func (b *BlockPool) removeElementAtStart() {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.size == 0 {
		return
	}
	if b.startIdx < 0 {
		return
	}

	k := b.startIdx
	if _, ok := b.pool[k]; ok {
		delete(b.pool, k)
	}

	if 0 == b.size {
		b.startIdx = -1
		b.endIdx = -1
	} else {
		b.startIdx++
	}
	b.size--
}
