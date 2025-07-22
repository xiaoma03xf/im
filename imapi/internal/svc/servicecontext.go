package svc

import (
	"zeroim/imapi/internal/config"
	"zeroim/imrpc/imrpc"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	IMRpc  imrpc.ImrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	userClient := zrpc.MustNewClient(c.ImRPC)
	return &ServiceContext{
		Config: c,
		IMRpc:  imrpc.NewImrpcClient(userClient.Conn()),
	}
}
