package shortlink

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"shortlink/data"
	"shortlink/pkg/gopool"
	"strconv"
	"sync"
	"time"
)

const (
	dicts    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	limitMax = 200
)

type LinkTask struct {
	lock     *RedisLock
	linkData *ShortLinkData
}

func NewLinkTask(lock *RedisLock, linkData *ShortLinkData) *LinkTask {
	return &LinkTask{
		lock:     lock,
		linkData: linkData,
	}
}

// StartRedisRecordDelete 启动redis zset过期短链删除任务
func (self *LinkTask) StartRedisRecordDelete(ctx context.Context, interval time.Duration) {
	timer := time.NewTimer(interval)
	for {
		select {
		case <-ctx.Done():
			slog.Info("redis过期短链删除任务终止")
			timer.Stop()
			return
		case <-timer.C:
			if ok, lockValue := self.lock.Lock(ctx, DeleteRedisLock, interval); ok {
				self.handlerRemoveRedisExpireRecord(ctx)
				if err := self.lock.Unlock(ctx, DeleteRedisLock, lockValue); err != nil {
					slog.Error("归还redis过期短链删除任务锁失败")
				}
			} else {
				slog.Info("获取redis过期短链删除任务锁失败，锁已被占用，结束处理！")
			}
			timer.Reset(interval)
		}
	}
}

// StartDbRecordDelete 启动数据库过期短链删除任务
func (self *LinkTask) StartDbRecordDelete(ctx context.Context, interval time.Duration) {
	timer := time.NewTimer(interval)
	for {
		select {
		case <-ctx.Done():
			slog.Info("数据库过期短链删除任务终止")
			timer.Stop()
			return
		case <-timer.C:
			if ok, lockValue := self.lock.Lock(ctx, DeleteDbLock, interval); ok {
				self.handlerRemoveDbExpireRecord(ctx)
				if err := self.lock.Unlock(ctx, DeleteDbLock, lockValue); err != nil {
					slog.Error("归还数据库过期短链删除任务锁失败")
				}
			} else {
				slog.Info("获取数据库过期短链删除任务锁失败，锁已被占用，结束处理！")
			}
			timer.Reset(interval)
		}
	}
}

func (self *LinkTask) handlerRemoveRedisExpireRecord(ctx context.Context) {
	slog.Info("开始执行redis过期短链删除任务")
	expireTime := time.Now().UnixMilli()
	var preLastScore float64
	for {
		var start any
		if preLastScore == 0 {
			start = "-inf"
		} else {
			start = preLastScore
		}
		// 通过偏移量和limit实现分批处理
		result, err := data.ReplicaRedis().ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
			Key:     data.ShortLinkExpireSet,
			Start:   start,
			Stop:    expireTime,
			Offset:  0,
			Count:   limitMax,
			ByScore: true,
		}).Result()
		if err != nil {
			slog.Warn("获取redis待删除的短链列表失败，跳过处理", slog.Any("error", err))
			break
		}
		size := len(result)
		if size == 0 {
			slog.Info("获取到redis待删除短链列表为空，结束当前处理任务")
			break
		}
		for i := 0; i < size; i++ {
			baseKey := result[i].Member.(string)
			// 异步向kafka中发送删除消息
			gopool.Go(func() {
				SendDeleteMessage(ctx, []byte(baseKey), nil)
			})
		}
		// 从过期列表中删除当前数据，避免在kafka消息还没消费前被其它节点的定时任务重复删除
		count, err := data.MasterRedis().ZRemRangeByScore(context.Background(), data.ShortLinkExpireSet,
			strconv.FormatFloat(result[0].Score, 'f', 0, 64),
			strconv.FormatFloat(result[size-1].Score, 'f', 0, 64)).Result()
		if err != nil {
			slog.Error("redis过期短链列表删除已处理元素失败", slog.Float64("minScore", result[0].Score),
				slog.Float64("maxScore", result[size-1].Score), slog.Any("error", err))
		} else {
			slog.Info("redis过期短链列表删除已处理元素成功", slog.Int64("deleteCount", count), slog.Int("size", size))
		}
		if size < limitMax {
			slog.Info("redis中待删除短链列表数据量小于limit限制，结束当前处理任务",
				slog.Int("size", size), slog.Int64("limit", limitMax))
		}
		preLastScore = result[size-1].Score
	}
	handleTime := time.Now().UnixMilli() - expireTime
	slog.Info("redis过期短链删除任务执行完成", slog.Int64("handlerTime", handleTime))
}

func (self *LinkTask) handlerRemoveDbExpireRecord(ctx context.Context) {
	slog.Info("开始执行数据库过期短链删除任务")
	expireTime := time.Now().UnixMilli()
	// 分批处理所有的表
	for i := 0; i < 62; i += 6 {
		wg := &sync.WaitGroup{}
		// 每张表启动一条协程进行删除
		for y := i; y < min(i+6, 62); y++ {
			wg.Add(1)
			suffix := string(dicts[y])
			gopool.Go(func() {
				defer wg.Done()
				self.removeSingleTableExpireRecord(ctx, expireTime, suffix)
			})
		}
		wg.Wait()
	}
	handleTime := time.Now().UnixMilli() - expireTime
	slog.Info("数据库过期短链删除任务执行完成", slog.Int64("handlerTime", handleTime))
}

func (self *LinkTask) removeSingleTableExpireRecord(ctx context.Context, expireTime int64, suffix string) {
	tableName := "\"\"\"" + TablePrefix + suffix + "\"\"\""
	var preLastId int64 = 0
	for {
		// 查询待删除的短链列表，通过游标分页进行查询，每次最多返回200条数据
		expireList, err := self.linkData.ListExpireLinkRecord(tableName, expireTime, preLastId, limitMax)
		if err != nil {
			slog.Warn("获取数据库待删除的短链列表失败", slog.String("tableName", tableName),
				slog.Int64("preLastId", preLastId), slog.Any("error", err))
			break
		}
		size := len(expireList)
		if size == 0 {
			slog.Info("获取到待删除短链列表为空，结束当前表的处理任务", slog.String("tableName", tableName))
			break
		}
		slog.Info("获取待删除的短链列表成功", slog.String("tableName", tableName),
			slog.Int64("preLastId", preLastId), slog.Int("count", len(expireList)))
		for _, record := range expireList {
			// 异步向kafka中发送删除短链消息
			gopool.Go(func() {
				SendDeleteMessage(ctx, []byte(record.BaseValue), nil)
			})
		}
		if size < limitMax {
			slog.Info("待删除短链列表数据量小于limit限制，结束当前表的处理任务", slog.String("tableName", tableName),
				slog.Int("size", size), slog.Int64("limit", limitMax))
			break
		}
		preLastId = expireList[size-1].Id
	}
}
