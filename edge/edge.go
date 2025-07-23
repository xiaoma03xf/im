package main

import (
	"context"
	"flag"
	"fmt"

	"zeroim/common/libnet"
	"zeroim/common/socket"
	"zeroim/edge/internal/config"
	"zeroim/edge/internal/logic"
	"zeroim/edge/internal/server"
	"zeroim/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	zeroservice "github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/edge.yaml", "the config file")

func main() {
	flag.Parse()

	var err error
	var c config.Config
	conf.MustLoad(*configFile, &c)
	srvCtx := svc.NewServiceContext(c)

	logx.DisableStat()

	tcpServer := server.NewTCPServer(srvCtx)
	protobuf := libnet.NewIMProtocol()
	tcpServer.Server, err = socket.NewServe(c.Name, c.TCPListenOn, protobuf, c.SendChanSize)
	if err != nil {
		panic(err)
	}

	go tcpServer.HandleRequest()
	go tcpServer.KqHeart()
	go tcpServer.RegisterEdge()

	fmt.Printf("Starting tcp server at %s \n", c.TCPListenOn)

	serviceGroup := zeroservice.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range logic.Consumers(context.Background(), srvCtx, tcpServer.Server) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
}
