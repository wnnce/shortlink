package balance

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"linkhttp/interceptor"
	"linkhttp/pb"
	"log/slog"
	"sync"
)

// GrpcClientPool GRPC连接池
type GrpcClientPool struct {
	connCache   map[string]*grpc.ClientConn     // 复用grpc连接
	clientCache map[string]pb.LinkServiceClient // 复用grpc client
	lock        sync.Mutex
}

func NewGrpcClientPool() *GrpcClientPool {
	return &GrpcClientPool{
		connCache:   make(map[string]*grpc.ClientConn),
		clientCache: make(map[string]pb.LinkServiceClient),
	}
}

func (self *GrpcClientPool) Client(addr string) (pb.LinkServiceClient, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if client, ok := self.clientCache[addr]; ok {
		return client, nil
	}
	conn, ok := self.connCache[addr]
	if !ok {
		newConn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(interceptor.TraceInterceptor))
		if err != nil {
			slog.Error("与rpc服务创建连接失败", slog.String("addr", addr), slog.Any("error", err))
			return nil, err
		}
		conn = newConn
		self.connCache[addr] = conn
	}
	client := pb.NewLinkServiceClient(conn)
	self.clientCache[addr] = client
	return client, nil
}

func (self *GrpcClientPool) Delete(addr string) {
	self.lock.Lock()
	delete(self.connCache, addr)
	delete(self.clientCache, addr)
	self.lock.Unlock()
}
