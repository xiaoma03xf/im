package logic

import (
	"context"
	"zeroim/common/libnet"
	"zeroim/common/session"
	"zeroim/common/socket"
	"zeroim/edge/internal/svc"
	"zeroim/imrpc/imrpc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"google.golang.org/protobuf/proto"
)

type MqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	server *socket.Server
	logx.Logger
}

func NewMqLogic(ctx context.Context, svcCtx *svc.ServiceContext, srv *socket.Server) *MqLogic {
	return &MqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		server: srv,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MqLogic) Consume(ctx context.Context, _, val string) error {
	var msg imrpc.PostMsg
	err := proto.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("[Consume] proto.Unmarshal val: %s error: %v", val, err)
		return err
	}
	logx.Infof("[Consume] succ msg: %+v body: %s", msg, msg.Msg)

	if len(msg.ToToken) > 0 {
		// ToToken, GetTokenSessions寻找当前kafka集群(eg: edge_01)
		// 消息发往当前kafka, kafka消费时应该确定把信息发往与其连接的哪个, tcp连接
		// Manager 管理当前节点下 edge_01 的所有连接，根据当前 token 获取对应的连接
		sessions := l.server.Manager.GetTokenSessions(msg.ToToken)
		for i := range sessions {
			s := sessions[i]
			if s == nil {
				logx.Errorf("[Consume] session not found, msg: %v\n", msg)
				continue
			}
			err := s.Send(makeMessage(&msg))
			if err != nil {
				logx.Errorf("[Consume] session send error, msg: %v, err: %v", msg, err)
			}
		}
	} else {
		s := l.server.Manager.GetSession(session.FromString(msg.SessionId))
		if s == nil {
			logx.Errorf("[Consume] session not found, msg: %v", msg)
			return nil
		}
		return s.Send(makeMessage(&msg))
	}

	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext, server *socket.Server) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConf, NewMqLogic(ctx, svcCtx, server)),
	}
}

func makeMessage(msg *imrpc.PostMsg) libnet.Message {
	var message libnet.Message
	message.Version = uint8(msg.Version)
	message.Status = uint8(msg.Status)
	message.ServiceId = uint16(msg.ServiceId)
	message.Cmd = uint16(msg.Cmd)
	message.Seq = msg.Seq
	message.Body = []byte(msg.Msg)
	return message
}
