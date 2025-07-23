package server

import (
	"context"
	"zeroim/common/discovery"
	"zeroim/common/socket"
	"zeroim/edge/client"
	"zeroim/edge/internal/svc"
	"zeroim/imrpc/imrpc"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type TCPServer struct {
	svcCtx *svc.ServiceContext
	Server *socket.Server
}

func NewTCPServer(svcCtx *svc.ServiceContext) *TCPServer {
	return &TCPServer{svcCtx: svcCtx}
}

func (srv *TCPServer) HandleRequest() {
	for {
		session, err := srv.Server.Accept()
		if err != nil {
			panic(err)
		}
		cli := client.NewClient(srv.Server.Manager, session, srv.svcCtx.IMRpc)
		go srv.sessionLoop(cli)
	}
}

func (srv *TCPServer) sessionLoop(client *client.Client) {
	message, err := client.Receive()
	if err != nil {
		logx.Errorf("[sessionLoop] client.Receive error: %v", err)
		_ = client.Close()
		return
	}

	// 登陆校验
	clientToken, sessionId, err := client.Login(message)
	if err != nil {
		logx.Errorf("[sessionLoop] client.Login error: %v", err)
		_ = client.Close()
		return
	}

	// 定时发送心跳给客户端保活
	go client.HeartBeat()
	srv.svcCtx.BizRedis.Hincrby(srv.Server.Name, "connections", 1)

	for {
		message, err := client.Receive()
		if err != nil {
			logx.Errorf("[sessionLoop] client.Receive error: %v", err)
			srv.svcCtx.IMRpc.Logout(context.Background(), &imrpc.LogoutRequest{
				Token:     clientToken,
				SessionId: sessionId,
			})
			srv.svcCtx.BizRedis.Hincrby(srv.Server.Name, "connections", -1)
			_ = client.Close()
			return
		}
		err = client.HandlePackage(message)
		if err != nil {
			logx.Errorf("[sessionLoop] client.HandleMessage error: %v", err)
		}
	}
}

func (srv *TCPServer) RegisterEdge() {
	// 所有edge 的地理信息写入 geo 集合下
	_, err := srv.svcCtx.BizRedis.GeoAdd("geo", &redis.GeoLocation{
		Name:      srv.Server.Name,
		Longitude: srv.svcCtx.Config.Longitude,
		Latitude:  srv.svcCtx.Config.Latitude,
	})
	if err != nil {
		logx.Errorf("[RegisterEdge] GeoAdd longitude:%v and latitude:%v err:%v", srv.svcCtx.Config.Longitude, srv.svcCtx.Config.Latitude, err)
		return
	}
	err = srv.svcCtx.BizRedis.Hset(srv.Server.Name, "address", srv.svcCtx.Config.TCPListenOn)
	if err != nil {
		logx.Errorf("[RegisterEdge] Hset address err:%v", err)
	}
}
func (srv *TCPServer) KqHeart() {
	work := discovery.NewQueueWorker(srv.svcCtx.Config.Etcd.Key, srv.svcCtx.Config.Etcd.Hosts, srv.svcCtx.Config.KqConf)
	work.HeartBeat()
}
