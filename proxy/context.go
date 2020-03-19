package proxy

import (
	"fmt"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/redis"
)

type route struct {
	ID       string
	Upstream string
}

type vault struct {
	ID string
}

func (v *vault) GetUpstream() string {
	return redis.Client().Get(fmt.Sprintf("vault:%s:upstream", v.ID)).Val()
}

func (v *vault) setUpstream(upstream string) {
	redis.Client().Set(fmt.Sprintf("vault:%s:upstream", v.ID), upstream, 0)
}

type userData struct {
	vault *vault
	route *route
}

func ctxUserData(ctx *goproxy.ProxyCtx) *userData {
	if ctx.UserData == nil {
		ctx.UserData = &userData{}
	}

	return ctx.UserData.(*userData)
}
