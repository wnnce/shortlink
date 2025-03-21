package server

import (
	"context"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"linkhttp/balance"
	"linkhttp/pb"
	"log/slog"
	"strconv"
	"time"
)

type HTTPServer struct {
	selectBalance balance.LoadBalance
	editBalance   balance.LoadBalance
}

func NewLinkServer(selectBalance, editBalance balance.LoadBalance) *HTTPServer {
	return &HTTPServer{
		selectBalance: selectBalance,
		editBalance:   editBalance,
	}
}

func (self *HTTPServer) ok(ctx *fasthttp.RequestCtx, body []byte) {
	ctx.Success("application/json", body)
}

func (self *HTTPServer) fail(ctx *fasthttp.RequestCtx, code int, message string) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
	ctx.SetBody(NewResult(code, message, nil))
}

func (self *HTTPServer) handlerRpcRequestError(ctx *fasthttp.RequestCtx, err error, defaultCode int, defaultMessage string) {
	rpcErr, ok := status.FromError(err)
	if !ok {
		self.fail(ctx, defaultCode, defaultMessage)
		return
	}
	code := defaultCode
	switch rpcErr.Code() {
	case codes.InvalidArgument, codes.OutOfRange:
		code = fasthttp.StatusBadRequest
		break
	case codes.NotFound, codes.Unimplemented:
		code = fasthttp.StatusNotFound
		break
	case codes.PermissionDenied:
		code = fasthttp.StatusUnauthorized
		break
	case codes.Internal, codes.DataLoss:
		code = fasthttp.StatusInternalServerError
	}
	self.fail(ctx, code, rpcErr.Message())
}

func (self *HTTPServer) makeDeadlineMetadata(deadTime time.Time, traceId uint64) (context.Context, context.CancelFunc) {
	c, cancel := context.WithDeadline(context.Background(), deadTime)
	md := metadata.Pairs("traceId", strconv.FormatUint(traceId, 10))
	return metadata.NewOutgoingContext(c, md), cancel
}

func (self *HTTPServer) ShortLinkRedirect(ctx *fasthttp.RequestCtx) {
	key := ctx.UserValue("key").(string)
	client, err := self.selectBalance.Get(key)
	if err != nil {
		self.fail(ctx, fasthttp.StatusInternalServerError, "服务连接失败，请重试")
		return
	}
	c, _ := self.makeDeadlineMetadata(time.Now().Add(3*time.Second), ctx.ID())
	request := pb.QueryRequest{}
	request.BaseKey = key
	resp, err := client.QueryOriginUrl(c, &request)
	if err != nil {
		slog.Error("rpc服务请求失败", slog.Uint64("traceId", ctx.ID()), slog.String("key", key), slog.Any("error", err))
		self.handlerRpcRequestError(ctx, err, fasthttp.StatusInternalServerError, "请求失败，请重试")
		return
	}
	ctx.Redirect(resp.OriginUrl, fasthttp.StatusMovedPermanently)
}

func (self *HTTPServer) CreateShortLink(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	if body == nil || len(body) == 0 {
		self.fail(ctx, fasthttp.StatusBadRequest, "请求数据不能为空")
		return
	}
	request := &pb.CreateRequest{}
	if err := json.Unmarshal(body, request); err != nil {
		self.fail(ctx, fasthttp.StatusBadRequest, "请求数据不合法")
		return
	}
	request.ClientIp = ctx.RemoteIP().String()
	request.UserAgent = string(ctx.Request.Header.UserAgent())
	client, err := self.editBalance.Get(balance.EmptyKey)
	if err != nil {
		self.fail(ctx, fasthttp.StatusInternalServerError, "服务连接失败，请重试")
		return
	}
	c, _ := self.makeDeadlineMetadata(time.Now().Add(5*time.Second), ctx.ID())
	resp, err := client.CreateShortLink(c, request)
	if err != nil {
		slog.Error("rpc服务创建短链失败", slog.Uint64("traceId", ctx.ID()), slog.Any("error", err))
		self.handlerRpcRequestError(ctx, err, fasthttp.StatusInternalServerError, "创建短链失败，请重试")
		return
	}
	self.ok(ctx, ResultByOk(resp.BaseKey))
}

func (self *HTTPServer) DeleteShortLink(ctx *fasthttp.RequestCtx) {
	key := ctx.UserValue("key").(string)
	client, err := self.editBalance.Get(balance.EmptyKey)
	if err != nil {
		self.fail(ctx, fasthttp.StatusInternalServerError, "服务连接失败，请重试")
		return
	}
	c, _ := self.makeDeadlineMetadata(time.Now().Add(5*time.Second), ctx.ID())
	request := &pb.DeleteRequest{
		BaseKey: key,
	}
	if _, err = client.DeleteShortLink(c, request); err != nil {
		slog.Error("rpc服务删除短链失败", slog.Uint64("traceId", ctx.ID()), slog.String("key", key), slog.Any("error", err))
		self.handlerRpcRequestError(ctx, err, fasthttp.StatusInternalServerError, "删除失败，请重试")
		return
	}
	self.ok(ctx, ResultByOk(nil))
}

func (self *HTTPServer) SelectShortLinkInfo(ctx *fasthttp.RequestCtx) {
	key := ctx.UserValue("key").(string)
	client, err := self.selectBalance.Get(key)
	if err != nil {
		self.fail(ctx, fasthttp.StatusInternalServerError, "服务连接失败，请重试")
		return
	}
	c, _ := self.makeDeadlineMetadata(time.Now().Add(5*time.Second), ctx.ID())
	request := &pb.QueryRequest{
		BaseKey: key,
	}
	record, err := client.QueryShortLinkInfo(c, request)
	if err != nil {
		slog.Error("rpc服务获取短链详情失败", slog.Uint64("traceId", ctx.ID()), slog.String("key", key), slog.Any("error", err))
		self.handlerRpcRequestError(ctx, err, fasthttp.StatusInternalServerError, "查询失败，请重试")
		return
	}
	self.ok(ctx, ResultByOk(record))
}
