package snowflake

import (
	"math/rand/v2"
	"shortlink/config"
	"sync"
	"time"
)

const (
	totalLength        = 64
	timestampLength    = 41
	AreaNumberLength   = 5
	DeviceNumberLength = 5
	serialNumberLength = 13
)

var (
	ConfigureReaderList = config.NewReaderList(ConfigureReader)
	defaultAreaId       = 0
	defaultNodeId       = 0
)

func ConfigureReader(bootstrap config.Bootstrap) (func(), error) {
	defaultAreaId = bootstrap.Server.AreaId
	defaultNodeId = bootstrap.Server.NodeId
	return nil, nil
}

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	areaId        int
	nodeId        int
	serial        uint32
	round         rand.Rand
}

func NewSnowflake(areaId, nodeId int) *Snowflake {
	return &Snowflake{
		areaId: areaId,
		nodeId: nodeId,
	}
}

func NewDefaultSnowflake() *Snowflake {
	return &Snowflake{
		areaId: defaultAreaId,
		nodeId: defaultNodeId,
	}
}

// GenerateId 通过雪花算法生成唯一的id
// 0-2固定为1
// 3-8位为随机数
// 9-41位为时间戳
// 42-46位为数据中心id
// 46-51位为机器id
// 51-63位为序列号
func (self *Snowflake) GenerateId() uint64 {
	self.mu.Lock()
	defer self.mu.Unlock()
	now := time.Now().UnixMilli()
	if now < self.lastTimestamp {
		self.lastTimestamp = now
	}
	if now == self.lastTimestamp && self.serial+1 < 8192 {
		self.serial++
	} else {
		self.serial = 0
		self.lastTimestamp = now
	}
	return 3<<62 |
		rand.Uint64N(63)<<56 |
		(uint64(now)<<31)>>8 |
		uint64(self.areaId)<<18 |
		uint64(self.nodeId)<<13 |
		uint64(self.serial)
}
