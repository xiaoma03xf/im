package logic

import (
	"context"
	"fmt"

	"zeroim/imapi/internal/svc"
	"zeroim/imapi/internal/types"
	"zeroim/imrpc/imrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMsgLogic) SendMsg(req *types.SendMsgRequest) (resp *types.SendMsgResponse, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.IMRpc.PostMessage(l.ctx, &imrpc.PostMsg{
		Token: fmt.Sprintf("%d", req.ToUserId),
		Msg:   req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &types.SendMsgResponse{}, nil
}
