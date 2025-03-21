package singleflight

import (
	"log"
	"sync"
)

type call struct {
	wg    sync.WaitGroup
	value any
	err   error
}

type Group struct {
	mu      sync.Mutex
	callMap map[string]*call
}

// Do 拦截多次key的相同请求
func (self *Group) Do(key string, fn func() (any, error)) (any, error) {
	self.mu.Lock()
	if self.callMap == nil {
		self.callMap = make(map[string]*call)
	}
	// 如果key已经在请求了，那么等待执行完成，直接返回
	if cl, ok := self.callMap[key]; ok {
		log.Println("peer is called key ", key)
		self.mu.Unlock()
		cl.wg.Wait()
		return cl.value, cl.err
	}
	newCall := &call{}
	newCall.wg.Add(1)
	self.callMap[key] = newCall
	self.mu.Unlock()
	newCall.value, newCall.err = fn()
	newCall.wg.Done()
	self.mu.Lock()
	delete(self.callMap, key)
	self.mu.Unlock()
	return newCall.value, newCall.err
}
