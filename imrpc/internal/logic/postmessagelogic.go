package logic

import (
	"context"
	"fmt"

	"zeroim/common/session"
	"zeroim/imrpc/imrpc"
	"zeroim/imrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type PostMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostMessageLogic {
	return &PostMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// func (l *PostMessageLogic) PostMessage(in *imrpc.PostMsg) (*imrpc.PostReponse, error) {

// 	// todo: add your logic here and delete this line
// 	ctx := context.Background()

// 	var (
// 		allDevice bool
// 		name      string
// 		token     string
// 		id        uint64
// 	)
// 	if len(in.Token) != 0 {
// 		allDevice = true
// 		token = in.Token
// 		logx.Infof("in.Token:%v, toToken:%v\n", in.Token, in.ToToken)
// 	} else {
// 		// sessionId = edge_name + ":" + token + ":" +id
// 		name, token, id = session.FromString(in.SessionId).Info()
// 	}
// 	// token -> [session1, session2, session3...],手机,电脑,平板多端登陆产生的sessionId
// 	sessionIds, err := l.svcCtx.BizRedis.Zrange(token, 0, -1)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(sessionIds) == 0 {
// 		return nil, err
// 	}
// 	for _, s := range sessionIds {
// 		name, token, id := session.FromString(s).Info()
// 		logx.Infof("sessionIds:%v, name:%v, token:%v, id:%v", s, name, token, id)
// 	}
// 	//把原始的 PostMsg 用 protobuf 编码成字节数组，准备投递
// 	data, err := proto.Marshal(in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// check data
// 	// logx.Infof("Version:%v, Status:%v\n", in.Version, in.Status)
// 	// logx.Infof("[PostMessage message to push: %v]", string(data))

// 	set := collection.NewSet()
// 	for _, sessionId := range sessionIds {
// 		respName, _, respId := session.FromString(sessionId).Info()
// 		// fmt.Println("respName:", respName, "respId", respId)
// 		if allDevice {
// 			set.Add(respName)
// 		} else {
// 			if name == respName && id == respId {
// 				edgeQueue, ok := l.svcCtx.QueueList.Load(respName)
// 				if !ok {
// 					logx.Severe("invalid session")
// 				} else {
// 					err = edgeQueue.Push(ctx, string(data))
// 					if err != nil {
// 						logx.Errorf("[PostMessage push data: %s error: %v]", string(data), err)
// 						return nil, err
// 					}
// 				}
// 			} else {
// 				logx.Severe("invalid session")
// 			}
// 		}
// 	}
// 	if set.Count() > 0 {
// 		logx.Infof("send to %d devices", set.Count())
// 	}

// 	// 广播
// 	for _, respName := range set.KeysStr() {
// 		// fmt.Println("查找edgeQueue:", respName)
// 		// l.svcCtx.QueueList.PrintlnKVS()
// 		edgeQueue, ok := l.svcCtx.QueueList.Load(respName)
// 		if !ok {
// 			logx.Errorf("invalid session")
// 		} else {
// 			err = edgeQueue.Push(ctx, string(data))
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}
// 	return &imrpc.PostReponse{}, nil
// }

func (l *PostMessageLogic) PostMessage(in *imrpc.PostMsg) (*imrpc.PostReponse, error) {

	// todo: add your logic here and delete this line
	if in.ToToken == "" {
		return nil, fmt.Errorf("[PostMessage] message toToken is nil")
	}
	sessions, err := l.svcCtx.BizRedis.Zrange(in.ToToken, 0, -1)
	if err != nil {
		return nil, err
	}
	data, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	set := collection.NewSet()
	for _, s := range sessions {
		toEdgeName, _, _ := session.FromString(s).Info()
		set.Add(toEdgeName)
	}

	if set.Count() > 0 {
		logx.Infof("send to %d devices", set.Count())
	}
	for _, respName := range set.KeysStr() {
		edgeQueue, ok := l.svcCtx.QueueList.Load(respName)
		if !ok {
			return nil, fmt.Errorf("[PostMessage] %v message pusher not found", respName)
		}
		err = edgeQueue.Push(context.Background(), string(data))
		if err != nil {
			logx.Errorf("[PostMessage] push data: %s error: %v", string(data), err)
		}
	}
	return &imrpc.PostReponse{}, nil
}
