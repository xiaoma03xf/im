package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	TCPListenOn  string
	SendChanSize int
	KqConf       kq.KqConf
	IMRpc        zrpc.RpcClientConf
}
