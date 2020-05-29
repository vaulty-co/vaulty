package proxy

import (
	"github.com/elazarl/goproxy"
	"github.com/vaulty/vaulty/routing"
)

type userData struct {
	route *routing.Route
}

func ctxUserData(ctx *goproxy.ProxyCtx) *userData {
	if ctx.UserData == nil {
		ctx.UserData = &userData{}
	}

	return ctx.UserData.(*userData)
}
