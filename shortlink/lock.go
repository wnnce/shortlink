package shortlink

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"shortlink/snowflake"
	"time"
)

const (
	DeleteRedisLock = "shortlink:task:redis:lock"
	DeleteDbLock    = "shortlink:task:db:lock"
)

type RedisLock struct {
	client *redis.Client
	snow   *snowflake.Snowflake
}

func NewRedisLock(client *redis.Client, snow *snowflake.Snowflake) *RedisLock {
	return &RedisLock{
		client: client,
		snow:   snow,
	}
}

func (self *RedisLock) Lock(ctx context.Context, key string, expireTime time.Duration) (bool, uint64) {
	value := self.snow.GenerateId()
	set, err := self.client.SetNX(ctx, key, value, expireTime).Result()
	if err != nil {
		slog.Error("redis分布式锁set异常", slog.String("key", key), slog.Uint64("value", value),
			slog.Any("error", err))
		return false, 0
	}
	return set, value
}

func (self *RedisLock) Unlock(ctx context.Context, key string, value uint64) error {
	script := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `
	result, err := self.client.Eval(ctx, script, []string{key}, value).Result()
	if err != nil {
		slog.Error("redis分布式锁del异常", slog.String("key", key), slog.Uint64("version", value),
			slog.Any("error", err))
		return err
	}
	if result == 0 {
		slog.Error("redis分布式锁解锁失败，锁已过期或version已更新", slog.String("key", key), slog.Uint64("value", value))
		return fmt.Errorf("锁过期或者value已被更新")
	}
	return nil
}
