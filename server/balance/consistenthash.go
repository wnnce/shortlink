package balance

import (
	"hash/crc32"
	"linkhttp/pb"
	"log/slog"
	"slices"
	"sort"
	"strconv"
	"sync"
)

type Hash func(key []byte) uint32

// HashLoadBalance 通过一致性哈希进行负载均衡
type HashLoadBalance struct {
	keys     []int          // 环形列表
	hashFunc Hash           // 哈希函数
	replicas int            // 节点的虚拟节点数量
	hashed   map[int]string // 节点映射
	lock     sync.RWMutex   // 读写锁
	pool     *GrpcClientPool
}

func NewHashLoadBalance(fn Hash, replicas int, pool *GrpcClientPool) LoadBalance {
	m := &HashLoadBalance{
		keys:     make([]int, 0),
		replicas: replicas,
		hashFunc: fn,
		hashed:   make(map[int]string),
		pool:     pool,
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// 元素在添加时，前面的元素已经是有序的了，所以使用插排算法效率最高
func (self *HashLoadBalance) sortKeys() {
	for i := 1; i < len(self.keys); i++ {
		base := self.keys[i]
		y := i - 1
		for y >= 0 && self.keys[y] > base {
			self.keys[y+1] = self.keys[y]
			y--
		}
		self.keys[y+1] = base
	}
}

// Add 添加节点主机
func (self *HashLoadBalance) Add(hosts ...HostInfo) {
	self.lock.Lock()
	defer self.lock.Unlock()
	for _, host := range hosts {
		for i := 1; i <= self.replicas; i++ {
			hashKey := int(self.hashFunc([]byte(host.Name + strconv.Itoa(i))))
			addr, ok := self.hashed[hashKey]
			if ok && addr == host.Addr {
				slog.Warn("远程主机已存在且地址相同，跳过添加", slog.String("hostname", host.Name), slog.String("address", addr))
				break
			}
			if ok {
				slog.Info("远程主机已存在但地址不同，更新主机地址", slog.String("hostname", host.Name), slog.String("address", addr))
				self.hashed[hashKey] = host.Addr
				continue
			}
			self.keys = append(self.keys, hashKey)
			self.hashed[hashKey] = host.Addr
		}
	}
	self.sortKeys()
}

func (self *HashLoadBalance) Get(key string) (pb.LinkServiceClient, error) {
	hashKey := int(self.hashFunc([]byte(key)))
	self.lock.RLock()
	index := sort.Search(len(self.keys), func(i int) bool {
		return self.keys[i] >= hashKey
	})
	hostAddr := self.hashed[self.keys[index%len(self.keys)]]
	self.lock.RUnlock()
	return self.pool.Client(hostAddr)
}

// Delete 节点主机下线时删除地址
func (self *HashLoadBalance) Delete(hostname string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	firstHashKey := int(self.hashFunc([]byte(hostname + "1")))
	addr, ok := self.hashed[firstHashKey]
	if !ok {
		return
	}
	firstIndex := sort.SearchInts(self.keys, firstHashKey)
	self.keys = slices.Delete(self.keys, firstIndex, firstIndex+1)
	delete(self.hashed, firstHashKey)
	for i := 2; i <= self.replicas; i++ {
		hashKey := int(self.hashFunc([]byte(hostname + strconv.Itoa(i))))
		delete(self.hashed, hashKey)
		hashKeyIndex := sort.SearchInts(self.keys, hashKey)
		self.keys = slices.Delete(self.keys, hashKeyIndex, hashKeyIndex+1)
	}
	self.pool.Delete(addr)

}

func (self *HashLoadBalance) Reset() {
	self.lock.Lock()
	self.keys = make([]int, 0)
	clear(self.hashed)
	self.lock.Unlock()
}
