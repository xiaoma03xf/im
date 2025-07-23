package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Longitude float64
	Latitude  float64

	Name         string
	TCPListenOn  string
	SendChanSize int

	IMRpc    zrpc.RpcClientConf
	KqConf   kq.KqConf
	Etcd     discov.EtcdConf
	BizRedis redis.RedisConf
}
