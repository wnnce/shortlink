package balance

import (
	"linkhttp/pb"
	"log/slog"
	"slices"
	"sort"
	"sync"
)

// PollLoadBalance 轮询负载均衡
type PollLoadBalance struct {
	hostnames []string
	index     int
	hostMap   map[string]string
	lock      sync.Mutex
	pool      *GrpcClientPool
}

func NewPollLoadBalance(pool *GrpcClientPool) LoadBalance {
	return &PollLoadBalance{
		hostnames: make([]string, 0),
		hostMap:   make(map[string]string),
		pool:      pool,
	}
}

func (self *PollLoadBalance) Add(hosts ...HostInfo) {
	self.lock.Lock()
	for _, host := range hosts {
		if addr, ok := self.hostMap[host.Name]; ok {
			if addr == host.Addr {
				slog.Warn("远程主机已存在且地址相同，跳过添加", slog.String("hostname", host.Name), slog.String("address", addr))
			} else {
				slog.Info("远程主机已存在但地址不同，更新主机地址", slog.String("hostname", host.Name), slog.String("address", addr))
				self.hostMap[host.Name] = host.Addr
			}
			continue
		}
		self.hostnames = append(self.hostnames, host.Name)
		self.hostMap[host.Name] = host.Addr
	}
	sort.Strings(self.hostnames)
	self.lock.Unlock()
}

func (self *PollLoadBalance) Get(_ string) (pb.LinkServiceClient, error) {
	self.lock.Lock()
	if self.index >= len(self.hostnames) {
		self.index = 0
	}
	curIndex := self.index
	self.index++
	hostAddr := self.hostMap[self.hostnames[curIndex]]
	self.lock.Unlock()
	return self.pool.Client(hostAddr)
}

func (self *PollLoadBalance) Delete(hostname string) {
	self.lock.Lock()
	if addr, ok := self.hostMap[hostname]; ok {
		index := sort.SearchStrings(self.hostnames, hostname)
		self.hostnames = slices.Delete(self.hostnames, index, index+1)
		delete(self.hostMap, hostname)
		self.pool.Delete(addr)
	}
	self.lock.Unlock()
}

func (self *PollLoadBalance) Reset() {
	self.lock.Lock()
	self.hostnames = make([]string, 0)
	clear(self.hostMap)
	self.index = 0
	self.lock.Unlock()
}
