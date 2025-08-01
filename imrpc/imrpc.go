package main

import (
	"flag"
	"fmt"

	"zeroim/imrpc/imrpc"
	"zeroim/imrpc/internal/config"
	"zeroim/imrpc/internal/server"
	"zeroim/imrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// goctl rpc protoc imrpc.proto --go_out=. --go-grpc_out=. --zrpc_out=.

var configFile = flag.String("f", "etc/imrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		imrpc.RegisterImrpcServer(grpcServer, server.NewImrpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
