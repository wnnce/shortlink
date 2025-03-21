package shortlink

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log/slog"
	"shortlink/config"
)

type ServiceRegister struct {
	client      *clientv3.Client
	leasId      clientv3.LeaseID
	keepChan    <-chan *clientv3.LeaseKeepAliveResponse
	serviceName string
	address     string
	ttl         int64
}

func NewServiceRegister(properties config.Etcd, serviceName, address string) (*ServiceRegister, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:            properties.Endpoints,
		DialTimeout:          properties.DialTimeout,
		DialKeepAliveTimeout: properties.DialKeepAliveTimeout,
		Username:             properties.Username,
		Password:             properties.Password,
	})
	if err != nil {
		return nil, err
	}
	return &ServiceRegister{
		client:      client,
		serviceName: serviceName,
		address:     address,
		ttl:         properties.TTL,
	}, nil
}

func (self *ServiceRegister) PutWithLease(ctx context.Context, leaseCallback func(resp *clientv3.LeaseKeepAliveResponse)) error {
	// 设置租约有效时间
	resp, err := self.client.Grant(ctx, self.ttl)
	if err != nil {
		return err
	}
	_, err = self.client.Put(ctx, self.serviceName, self.address, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	slog.Info("服务注册成功", slog.String("name", self.serviceName), slog.String("address", self.address),
		slog.Any("leaseId", resp.ID))
	keepChan, err := self.client.KeepAlive(ctx, resp.ID)
	if err != nil {
		return err
	}
	self.leasId = resp.ID
	self.keepChan = keepChan
	if leaseCallback != nil {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case leaseResp := <-keepChan:
					leaseCallback(leaseResp)
				}
			}
		}()
	}
	return nil
}

func (self *ServiceRegister) Close() {
	self.client.Revoke(context.Background(), self.leasId)
	self.client.Close()
}
