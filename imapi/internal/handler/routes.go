// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.5

package handler

import (
	"net/http"

	"zeroim/imapi/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/send/msg",
				Handler: SendMsgHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1/imapi"),
	)
}
