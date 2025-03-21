package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"hash/crc32"
	"linkhttp"
	"linkhttp/balance"
	"linkhttp/config"
	"log"
	"log/slog"
)

var configFilePath string

func main() {
	flag.StringVar(&configFilePath, "config", "/config/config-dev.yaml", "server config file path")
	flag.Parse()

	if err := config.InitConfig(configFilePath); err != nil {
		panic(err)
	}
	pool := balance.NewGrpcClientPool()
	pollBalance := balance.NewPollLoadBalance(pool)
	hashBalance := balance.NewHashLoadBalance(crc32.ChecksumIEEE, 8, pool)
	// 注册服务发现
	disc, err := server.NewServiceDiscover(config.DefaultConfig().Etcd)
	if err != nil {
		panic(err)
	}
	serverConfig := config.DefaultConfig().Server
	ctx, cancel := context.WithCancel(context.Background())
	if err = disc.WatchService(ctx, serverConfig.ServiceName, func(name, addr string) {
		host := balance.HostInfo{Name: name, Addr: addr}
		pollBalance.Add(host)
		hashBalance.Add(host)
	}, func(name string) {
		pollBalance.Delete(name)
		hashBalance.Delete(name)
	}); err != nil {
		panic(err)
	}
	linkServer := server.NewLinkServer(hashBalance, pollBalance)

	// 路由绑定
	rou := router.New()
	rou.POST("/link", linkServer.CreateShortLink)
	rou.DELETE("/link/{key:[0-9A-Za-z]{11}}", linkServer.DeleteShortLink)
	rou.GET("/info/{key:[0-9A-Za-z]{11}}", linkServer.SelectShortLinkInfo)
	rou.GET("/{key:[0-9A-Za-z]{11}}", linkServer.ShortLinkRedirect)

	defer func() {
		cancel()
	}()
	address := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	slog.Info("http服务启动 listen " + address)
	log.Fatalln(fasthttp.ListenAndServe(address, rou.Handler))
}
