package shortlink

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"shortlink/pb"
)

// ServiceServer rpc服务
type ServiceServer struct {
	pb.UnimplementedLinkServiceServer
	service *ShortLinkService
}

func NewServiceServer(service *ShortLinkService) pb.LinkServiceServer {
	return &ServiceServer{
		service: service,
	}
}

func (self *ServiceServer) QueryOriginUrl(ctx context.Context, request *pb.QueryRequest) (*pb.QueryResponse, error) {
	originUrl, err := self.service.SelectOriginUrlByKey(ctx, request.BaseKey)
	if err != nil {
		return nil, err
	}
	return &pb.QueryResponse{
		OriginUrl: originUrl,
	}, nil
}

func (self *ServiceServer) CreateShortLink(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	shortLink, err := self.service.PreAddShortLink(ctx, request)
	if err != nil {
		return nil, err
	}
	return &pb.CreateResponse{
		BaseKey: shortLink,
	}, nil
}

func (self *ServiceServer) DeleteShortLink(ctx context.Context, request *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := self.service.DeleteShortLinkByKey(ctx, request.BaseKey)
	return &emptypb.Empty{}, err
}

func (self *ServiceServer) QueryShortLinkInfo(ctx context.Context, request *pb.QueryRequest) (*pb.LinkRecord, error) {
	return self.service.SelectInfoByKey(ctx, request.BaseKey)
}
