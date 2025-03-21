package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

// TraceInterceptor grpc日志记录拦截器
func TraceInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	begin := time.Now().UnixMilli()
	err := invoker(ctx, method, req, reply, cc, opts...)
	handleTime := time.Now().UnixMilli() - begin
	slog.Info(fmt.Sprintf("[ %s ] -> [ method: %s ] %dms 处理完成", cc.Target(), method, handleTime))
	return err
}
