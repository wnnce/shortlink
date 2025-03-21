package server

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"linkhttp/config"
	"log/slog"
	"sync"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	client *clientv3.Client
	lock   sync.Mutex
}

func NewServiceDiscover(properties config.Etcd) (*ServiceDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   properties.Endpoints,
		DialTimeout: properties.DialTimeout,
		Username:    properties.Username,
		Password:    properties.Password,
	})
	if err != nil {
		return nil, err
	}
	return &ServiceDiscovery{
		client: client,
	}, nil
}

// 服务启动时先获取当前在线的所有rpc节点
func (self *ServiceDiscovery) initServiceAddress(ctx context.Context, prefix string, callback func(name, add string)) error {
	resp, err := self.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		callback(string(kv.Key), string(kv.Value))
	}
	slog.Info("etcd服务注册成功", slog.Any("hostCount", resp.Count))
	return nil
}

// WatchService 监听rpc节点上线或下线
func (self *ServiceDiscovery) WatchService(ctx context.Context, prefix string, putCallback func(name, addr string),
	delCallback func(name string)) error {
	if err := self.initServiceAddress(ctx, prefix, putCallback); err != nil {
		return err
	}
	watchChan := self.client.Watch(ctx, prefix, clientv3.WithPrefix())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case watchResp := <-watchChan:
				for _, ev := range watchResp.Events {
					key, value := string(ev.Kv.Key), string(ev.Kv.Value)
					switch ev.Type {
					case mvccpb.PUT:
						slog.Info("etcd通知服务上线", slog.String("serviceName", key), slog.String("address", value))
						go putCallback(key, value)
					case mvccpb.DELETE:
						slog.Info("etcd通知服务下线", slog.String("serviceNam", key), slog.String("address", value))
						go delCallback(key)
					}
				}
			}
		}
	}()
	return nil
}
