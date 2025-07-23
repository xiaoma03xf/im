package svc

import (
	"zeroim/edge/internal/config"
	"zeroim/imrpc/imrpcclient"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	IMRpc    imrpcclient.Imrpc
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.IMRpc)
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:   c,
		IMRpc:    imrpcclient.NewImrpc(client),
		BizRedis: rds,
	}
}
