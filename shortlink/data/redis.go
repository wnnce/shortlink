package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"shortlink/config"
	"time"
)

const (
	shortLinkPrefix    = "shortlink:"
	ShortLinkExpireSet = "shortlink:expire:set"
	ShortLinkFilter    = "shortlink:filter"
)

func RedisConfigureReader(bootstrap config.Bootstrap) (func(), error) {
	masterClient, err := makeRedisClient(bootstrap.Data.Redis.Master)
	if err != nil {
		return nil, err
	}
	clients := make([]*redis.Client, 0, len(bootstrap.Data.Redis.Replicas))
	for _, replica := range bootstrap.Data.Redis.Replicas {
		replicaClient, err := makeRedisClient(replica)
		if err != nil {
			return nil, err
		}
		clients = append(clients, replicaClient)
	}
	defaultRedisData = &redisData{
		master:      masterClient,
		replicas:    clients,
		loadBalance: NewPollLoadBalance(int32(len(clients))),
	}
	return func() {
		masterClient.Close()
		for _, client := range clients {
			client.Close()
		}
	}, nil
}

func makeRedisClient(redisConfig config.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		DB:           redisConfig.Index,
		Username:     redisConfig.Username,
		Password:     redisConfig.Password,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("redis服务端通信失败", slog.String("host", redisConfig.Host), slog.Int("port", redisConfig.Port),
			slog.String("username", redisConfig.Username), slog.String("password", redisConfig.Password),
			slog.Int("dbIndex", redisConfig.Index))
		return nil, err
	}
	slog.Info("redis客户端创建成功", slog.String("host", redisConfig.Host), slog.Int("port", redisConfig.Port),
		slog.String("username", redisConfig.Username), slog.String("password", redisConfig.Password),
		slog.Int("dbIndex", redisConfig.Index))
	return rdb, nil
}

var defaultRedisData *redisData

type redisData struct {
	master      *redis.Client
	replicas    []*redis.Client
	loadBalance LoadBalance
}

func (self *redisData) MasterClient() *redis.Client {
	return self.master
}

func (self *redisData) ReplicaClients() []*redis.Client {
	return self.replicas
}

func (self *redisData) ReplicaClient() *redis.Client {
	if self.replicas == nil || len(self.replicas) == 0 {
		return self.master
	}
	index := self.loadBalance.GetIndex()
	return self.replicas[index]
}

func MasterRedis() *redis.Client {
	if defaultRedisData == nil {
		return nil
	}
	return defaultRedisData.MasterClient()
}

func ReplicaRedis() *redis.Client {
	if defaultRedisData == nil {
		return nil
	}
	return defaultRedisData.ReplicaClient()
}

// RedisGetStruct 使用传递的client查询在Redis中缓存的结构体
// 使用泛型指定返回结构体类型
func RedisGetStruct[T any](ctx context.Context, key string, client *redis.Client) (T, error) {
	value := new(T)
	result, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return *value, err
	}
	err = json.Unmarshal(result, value)
	return *value, err
}

// RedisGetSlice 使用默认redisTemplate查询在redis中缓存的切片
// 如果默认redisTemplate为nil 那么会报nil空地址异常
// 使用泛型指定切片的类型 泛型可以为指针
func RedisGetSlice[T any](ctx context.Context, key string, client *redis.Client) ([]T, error) {
	result, err := client.Get(ctx, key).Bytes()
	if err != nil {
		slog.Error("查询Redis缓存结构体失败", "error", err.Error(), "key", key)
		return nil, err
	}
	value := make([]T, 0)
	err = json.Unmarshal(result, &value)
	return value, err
}

// GetShortLinkKey 获取指定链接key的redis缓存key
func GetShortLinkKey(key string) string {
	return shortLinkPrefix + key
}
