package main

import (
	"context"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"shortlink"
	"shortlink/config"
	"shortlink/data"
	"shortlink/pb"
	"shortlink/snowflake"
	"strconv"
	"time"
)

func HandlerEtcdLease(resp *clientv3.LeaseKeepAliveResponse) {
	slog.Debug("etcd租期续约成功", slog.Any("leaseId", resp.ID))
}

var configFilePath string

func main() {
	flag.StringVar(&configFilePath, "config", "/config/config-dev.yaml", "server config file path")
	flag.Parse()
	if err := config.InitConfig(configFilePath); err != nil {
		panic(err)
	}
	// 分发配置
	config.RegisterConfigureReader(data.ConfigureReaderList, snowflake.ConfigureReaderList, shortlink.ConfigureReaderList)
	cleanup := config.IssueConfigure()
	ctx, cancel := context.WithCancel(context.Background())

	serverConfig := config.DefaultConfig().Server
	address := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	// 向etcd注册服务
	register, err := shortlink.NewServiceRegister(config.DefaultConfig().Etcd,
		serverConfig.ServiceName+"-"+strconv.Itoa(serverConfig.NodeId), address)
	if err != nil {
		panic(err)
	}
	if err = register.PutWithLease(ctx, HandlerEtcdLease); err != nil {
		panic(err)
	}

	// 注册rpc服务
	server := grpc.NewServer()
	linkData := shortlink.NewShortLinkData()
	linkService := shortlink.NewShortLinkService(snowflake.NewDefaultSnowflake(), linkData)
	linkServer := shortlink.NewServiceServer(linkService)
	linkTask := shortlink.NewLinkTask(shortlink.NewRedisLock(data.MasterRedis(), snowflake.NewDefaultSnowflake()), linkData)
	pb.RegisterLinkServiceServer(server, linkServer)

	// 监听kafka消息
	go shortlink.StartAddConsumer(ctx, linkService)
	go shortlink.StartDeleteConsumer(ctx, linkService)

	// 启动短链删除任务
	go linkTask.StartRedisRecordDelete(ctx, 15*time.Minute)
	go linkTask.StartDbRecordDelete(ctx, 6*time.Hour)

	defer func() {
		cancel()
		cleanup()
		register.Close()
	}()

	// 启动服务
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
	if err != nil {
		panic(err)
	}
	log.Fatalln(server.Serve(listen))
}
