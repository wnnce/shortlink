package shortlink

import (
	"context"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"shortlink/config"
	"shortlink/pb"
)

var addWrite, deleteWrite *kafka.Writer
var addReader, deleteReader *kafka.Reader

func KafkaConfigureReader(bootstrap config.Bootstrap) (func(), error) {
	addWrite = &kafka.Writer{
		Addr:                   kafka.TCP(bootstrap.Kafka.Brokers...),
		Topic:                  bootstrap.Kafka.AddTopic,
		AllowAutoTopicCreation: true,
	}
	deleteWrite = &kafka.Writer{
		Addr:                   kafka.TCP(bootstrap.Kafka.Brokers...),
		Topic:                  bootstrap.Kafka.DeleteTopic,
		AllowAutoTopicCreation: true,
	}
	addReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  bootstrap.Kafka.Brokers,
		GroupID:  bootstrap.Kafka.AddGroupId,
		Topic:    bootstrap.Kafka.AddTopic,
		MaxBytes: 10e6,
	})
	deleteReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  bootstrap.Kafka.Brokers,
		GroupID:  bootstrap.Kafka.DeleteGroupId,
		Topic:    bootstrap.Kafka.DeleteTopic,
		MaxBytes: 10e6,
	})
	return func() {
		addWrite.Close()
		deleteWrite.Close()
		addReader.Close()
		deleteReader.Close()
	}, nil
}

// StartAddConsumer 启动新增链接消费者
func StartAddConsumer(ctx context.Context, service *ShortLinkService) {
	var fetchErrorCount int64
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := addReader.ReadMessage(ctx)
			if err != nil {
				slog.Error("kafka读取消息失败", slog.Any("error", err), slog.Int64("tryAgain", fetchErrorCount))
				if fetchErrorCount == 3 {
					slog.Error("kafka读取消息连续失败超过3次，终止读取")
					return
				}
				fetchErrorCount++
				continue
			}
			fetchErrorCount = 0
			slog.Info("接收到kafka链接新增消息", slog.String("key", string(message.Key)),
				slog.Int64("offset", message.Offset), slog.Int("partition", message.Partition))
			linkRecord := &pb.LinkRecord{}
			if err = proto.Unmarshal(message.Value, linkRecord); err != nil {
				slog.Error("kafka消息反序列化失败，删除redis缓存", slog.String("key", string(message.Key)),
					slog.String("value", string(message.Value)), slog.Any("error", err))
				service.ClearShortLinkCacheAndFilter(string(message.Key))
				continue
			}
			// 保存到数据库中
			service.AddShortLinkToDb(context.Background(), linkRecord)
		}
	}
}

// StartDeleteConsumer 启动删除链接消费者
func StartDeleteConsumer(ctx context.Context, service *ShortLinkService) {
	var fetchErrorCount int64
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := deleteReader.FetchMessage(ctx)
			if err != nil {
				slog.Error("kafka读取消息失败", slog.Any("error", err), slog.Int64("tryAgain", fetchErrorCount))
				if fetchErrorCount == 3 {
					slog.Error("kafka读取消息连续失败超过3次，终止读取")
					return
				}
				fetchErrorCount++
				continue
			}
			fetchErrorCount = 0
			baseKey := string(message.Key)
			slog.Info("接收到kafka链接删除消息", slog.String("key", baseKey),
				slog.Int64("offset", message.Offset), slog.Int("partition", message.Partition))
			// 删除失败不提交偏移量
			if err = service.DeleteShortLinkByMessageKey(baseKey); err != nil {
				if err = deleteReader.CommitMessages(context.Background(), message); err != nil {
					slog.Error("kafka提交链接删除消息偏移量失败", slog.String("key", baseKey),
						slog.Int64("offset", message.Offset), slog.Any("error", err))
				}
			}
		}
	}
}

// SendAddMessage 向kafka发送新增链接消息
// callback 发送失败的回调函数
func SendAddMessage(ctx context.Context, key, body []byte, callback func()) {
	if err := addWrite.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: body,
	}); err != nil {
		slog.Error("kafka发送新增消息失败", slog.String("key", string(key)), slog.Any("error", err))
		if callback != nil {
			callback()
		}
	}
}

// SendDeleteMessage 向kafka发送删除短链消息
func SendDeleteMessage(ctx context.Context, key []byte, callback func()) {
	if err := deleteWrite.WriteMessages(ctx, kafka.Message{
		Key: key,
	}); err != nil {
		slog.Error("kafka发送删除消息失败", slog.String("key", string(key)), slog.Any("error", err))
		if callback != nil {
			callback()
		}
	}
}
