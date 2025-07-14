package main

import (
	"flag"
	"fmt"

	"zeroim/edge/edge"
	"zeroim/edge/internal/config"
	"zeroim/edge/internal/server"
	"zeroim/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/edge.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		edge.RegisterEdgeServer(grpcServer, server.NewEdgeServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	tcpServer, err := server.NewTCPServer(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		tcpServer.Close()
	}()
	fmt.Printf("Starting tcp server at %s...\n", c.TCPListenOn)
	threading.GoSafe(func() {
		tcpServer.HandleRequest()
	})

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
