package shortlink

import (
	"context"
	"fmt"
	"github.com/jxskiss/base62"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"shortlink/data"
	"shortlink/pb"
	"shortlink/pkg/gopool"
	"shortlink/singleflight"
	"shortlink/snowflake"
	"strconv"
	"time"
)

type ShortLinkService struct {
	snowId    *snowflake.Snowflake
	linkRepo  *ShortLinkData
	callGroup *singleflight.Group
}

func NewShortLinkService(snowId *snowflake.Snowflake, linkRepo *ShortLinkData) *ShortLinkService {
	return &ShortLinkService{
		snowId:    snowId,
		linkRepo:  linkRepo,
		callGroup: &singleflight.Group{},
	}
}

// ClearShortLinkCacheAndFilter 清除预生成短链在redis和cuckoo过滤器中的数据
func (self *ShortLinkService) ClearShortLinkCacheAndFilter(key string) {
	// 从缓存中删除
	data.MasterRedis().Del(context.Background(), data.GetShortLinkKey(key))
	// 从cuckoo过滤器中删除
	data.MasterRedis().CFDel(context.Background(), data.ShortLinkFilter, key)
}

// PreAddShortLink 通过请求参数预生成短链
// 返回生成后的短链和rpc错误
func (self *ShortLinkService) PreAddShortLink(ctx context.Context, form *pb.CreateRequest) (string, error) {
	uniqueId := self.snowId.GenerateId()
	baseKey := base62.FormatUint(uniqueId)
	stringKey := string(baseKey)
	// 将生成的短链保存到缓存
	if _, err := data.MasterRedis().Set(ctx, data.GetShortLinkKey(stringKey), form.OriginUrl, 0).Result(); err != nil {
		slog.Error("添加短链redis缓存失败", slog.String("key", stringKey), slog.String("originUrl", form.OriginUrl),
			slog.Any("error", err))
		return "", status.Error(codes.Internal, "短链生成失败，请重试！")
	}
	// 保存到cuckoo过滤器
	if _, err := data.MasterRedis().CFAdd(ctx, data.ShortLinkFilter, stringKey).Result(); err != nil {
		slog.Error("短链添加到cuckoo过滤器失败", slog.String("key", stringKey), slog.Any("error", err))
		data.MasterRedis().Del(ctx, data.GetShortLinkKey(stringKey))
		return "", status.Error(codes.Internal, "短链生成失败，请重试！")
	}
	now := time.Now()
	linkRecord := &pb.LinkRecord{
		UniqueId:   strconv.FormatUint(uniqueId, 10),
		BaseValue:  stringKey,
		OriginUrl:  form.OriginUrl,
		IsLasting:  form.IsLasting,
		CreateTime: now.Format("2006-01-02 15:04:05"),
		ClientIp:   form.ClientIp,
		UserAgent:  form.UserAgent,
	}
	if !linkRecord.IsLasting {
		expireTime := now.Add(time.Duration(form.ValidHour) * time.Hour).UnixMilli()
		linkRecord.ValidHour = form.ValidHour
		linkRecord.ExpireTime = expireTime
		linkRecord.ExpireMode = form.ExpireMode
	}
	recordBytes, _ := proto.Marshal(linkRecord)
	// 将新增的短链记录异步发送到kafka
	gopool.Go(func() {
		SendAddMessage(context.Background(), baseKey, recordBytes, func() {
			// 消息发送失败的回调，如果失败那么删除redis和cuckoo过滤器中的缓存数据
			self.ClearShortLinkCacheAndFilter(string(baseKey))
		})
	})
	return stringKey, nil
}

// AddShortLinkToDb 将预生成的短链保存到数据库
func (self *ShortLinkService) AddShortLinkToDb(_ context.Context, record *pb.LinkRecord) {
	if err := self.linkRepo.Add(record); err != nil {
		slog.Error("链接记录保存到数据库失败", slog.String("key", record.BaseValue),
			slog.String("originUrl", record.OriginUrl), slog.Any("error", err))
		self.ClearShortLinkCacheAndFilter(record.BaseValue)
		return
	}
	slog.Info("链接记录保存到数据库成功", slog.Int64("id", record.Id),
		slog.String("key", record.BaseValue))
	// 将非长期链接且为精确删除模式的短链添加到redis的过期列表中
	if !record.IsLasting && record.ExpireMode == 1 {
		if _, err := data.MasterRedis().ZAdd(context.Background(), data.ShortLinkExpireSet, redis.Z{
			Score:  float64(record.ExpireTime),
			Member: record.BaseValue,
		}).Result(); err != nil {
			// 添加失败只输出日志，还有数据库的定时扫描删除做兜底处理
			slog.Info("短链记录添加到redis过期列表中失败", slog.String("key",
				record.BaseValue), slog.Any("error", err))
		}
	}
}

func (self *ShortLinkService) SelectOriginUrlByKey(ctx context.Context, key string) (string, error) {
	// 先查询cuckoo过滤器
	if ok, err := data.ReplicaRedis().CFExists(ctx, data.ShortLinkFilter, key).Result(); err != nil {
		slog.Error("查询cuckoo过滤器失败", slog.String("key", key), slog.Any("error", err))
	} else if !ok {
		slog.Info("短链在redis过滤器中不存在，跳过获取", slog.String("key", key))
		return "", status.Error(codes.NotFound, "此短链不存在！")
	}
	// 再查询缓存
	if result, err := data.ReplicaRedis().Get(ctx, data.GetShortLinkKey(key)).Result(); err == nil {
		// 如果cuckoo过滤器中存在，但是在缓存和数据库中都没有查询到，那么缓存一个默认的空值，避免直接请求数据库
		if result == "not" {
			return "", status.Error(codes.NotFound, "此短链不存在！")
		}
		return result, nil
	} else {
		slog.Error("查询redis短链缓存失败，开始查询数据库", slog.String("key", key), slog.Any("error", err))
	}
	// 使用call Do进行加锁查询，确保某个key在被多个请求同时查询时只有一个请求会查询数据库，其余请求均重用这次的数据库查询结果
	findResult, err := self.callGroup.Do(key, func() (any, error) {
		record, dbErr := self.linkRepo.SelectByKey(key)
		if dbErr != nil {
			slog.Error("查询数据库短链记录失败", slog.String("key", key), slog.Any("error", dbErr))
		}
		// 如果记录不存在，向redis中添加一个短时间过期的默认值
		if record == nil {
			data.MasterRedis().Set(context.Background(), data.GetShortLinkKey(key), "not", 1*time.Minute)
		} else {
			data.MasterRedis().Set(context.Background(), key, record.OriginUrl, 0)
		}
		return record, dbErr
	})
	if err != nil {
		return "", status.Error(codes.Internal, "查询失败，请重试！")
	}
	if findResult == nil {
		return "", status.Error(codes.NotFound, "此短链不存在！")
	}
	// 进行类型转换
	record, ok := findResult.(*pb.LinkRecord)
	if record == nil || !ok {
		return "", status.Error(codes.NotFound, "此短链不存在！")
	}
	return record.OriginUrl, nil
}

func (self *ShortLinkService) SelectInfoByKey(ctx context.Context, key string) (*pb.LinkRecord, error) {
	ok, err := data.ReplicaRedis().CFExists(ctx, data.ShortLinkFilter, key).Result()
	if err != nil {
		slog.Error("获取短链详情查询cuckoo过滤器失败", slog.String("key", key), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "查询失败，请重试")
	}
	if !ok {
		slog.Warn("cuckoo过滤器中此短链不存在，跳过查询", slog.String("key", key))
		return nil, status.Errorf(codes.NotFound, "该短链不存在")
	}
	// 查询短链详情仅为管理人员可用，不考虑缓存
	record, err := self.linkRepo.SelectInfoByKey(key)
	if err != nil {
		slog.Error("数据库查询短链记录失败", slog.String("key", key), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "查询失败，请重试")
	}
	return record, nil
}

// DeleteShortLinkByMessageKey 处理短链过期删除任务发送的key
func (self *ShortLinkService) DeleteShortLinkByMessageKey(key string) error {
	// 先删除redis缓存，如果数据库删除失败也不会影响最终的查询
	if _, err := data.MasterRedis().Del(context.Background(), data.GetShortLinkKey(key)).Result(); err != nil {
		slog.Error("redis删除链接key失败", slog.String("key", key), slog.Any("error", err))
		return fmt.Errorf("缓存删除短链失败")
	}
	// 再删除数据库中的记录
	if _, err := self.linkRepo.DeleteByBaseKey(key); err != nil {
		slog.Error("数据库删除链接key失败", slog.String("key", key), slog.Any("error", err))
		return fmt.Errorf("数据库删除短链失败")
	}
	// 删除布谷过滤器中的记录, 如果删除失败也不影响，因为数据库和缓存中已经删除
	if _, err := data.MasterRedis().CFDel(context.Background(), data.ShortLinkFilter, key).Result(); err != nil {
		slog.Error("redis过滤器删除key失败", slog.String("key", key), slog.Any("error", err))
	}
	return nil
}

// DeleteShortLinkByKey 处理短链删除请求
func (self *ShortLinkService) DeleteShortLinkByKey(ctx context.Context, key string) error {
	if ok, err := data.MasterRedis().CFExists(ctx, data.ShortLinkFilter, key).Result(); err != nil {
		slog.Error("查询cuckoo过滤器失败", slog.String("key", key), slog.Any("error", err))
	} else if !ok {
		slog.Info("请求删除的短链key不存在，跳过处理", slog.String("key", key))
		return status.Errorf(codes.InvalidArgument, "请求删除的短链不存在")
	}
	if err := self.DeleteShortLinkByMessageKey(key); err != nil {
		return status.Errorf(codes.Internal, "删除失败，请重试")
	}
	// 短链可能会在redis的短链过期列表中，也尝试删除，删除失败也不影响，后续被过期任务读取后会再进行删除
	// 同时删除是幂等操作，重复删除不会发生数据一致性问题
	data.MasterRedis().ZRem(context.Background(), data.ShortLinkExpireSet, key)
	return nil
}
