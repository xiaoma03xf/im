package server

import (
	"fmt"
	"zeroim/common/discovery"
	"zeroim/common/libnet"
	"zeroim/common/socket"
	"zeroim/edge/internal/svc"
)

type TCPServer struct {
	svcCtx *svc.ServiceContext
	Server *socket.Server
}

func NewTCPServer(svcCtx *svc.ServiceContext) (*TCPServer, error) {
	protocol := libnet.NewIMProtocol()
	s, err := socket.NewServe(svcCtx.Config.Name, svcCtx.Config.TCPListenOn, protocol, svcCtx.Config.SendChanSize)
	if err != nil {
		return nil, err
	}
	return &TCPServer{svcCtx: svcCtx, Server: s}, nil
}

func (srv *TCPServer) HandleRequest() {
	for {
		session, err := srv.Server.Accept()
		if err != nil {
			panic(err)
		}
		go srv.handleRequest(session)
	}
}

func (srv *TCPServer) handleRequest(session *libnet.Session) {
	defer session.Close()
	for {
		msg, err := session.Receive()
		if err != nil {
			fmt.Println("Receive err")
			panic(err)
		}
		// TODO 直接返回一样的内容
		if err := session.SendSample(*msg); err != nil {
			fmt.Println("resend err")
			panic(err)
		}
	}
}
func (srv *TCPServer) Close() {
	srv.Server.Close()
}

func (srv *TCPServer) KqHeart() {
	work := discovery.NewQueueWorker(srv.svcCtx.Config.Etcd.Key, srv.svcCtx.Config.Etcd.Hosts, srv.svcCtx.Config.KqConf)
	work.HeartBeat()
}
