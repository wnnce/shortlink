package data

import "sync/atomic"

// LoadBalance 数据库查询负载均衡器
type LoadBalance interface {
	GetIndex() int
}

// PollLoadBalance 简单的轮询负载均衡器
type PollLoadBalance struct {
	cap   int32
	index int32
}

func NewPollLoadBalance(cap int32) LoadBalance {
	return &PollLoadBalance{
		cap:   cap,
		index: 0,
	}
}

func (self *PollLoadBalance) GetIndex() int {
	currentIndex := atomic.LoadInt32(&self.index)
	nextIndex := (currentIndex + 1) % self.cap
	atomic.StoreInt32(&self.index, nextIndex)
	return int(currentIndex)
}
