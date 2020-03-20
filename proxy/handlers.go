package proxy

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
)

func (p *Proxy) NotFound(message string) goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return nil, errResponse(req, message, http.StatusNotFound)
	})
}

func (p *Proxy) TransformRequestBody() goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		fmt.Println("route found: ", ctxUserData(ctx).route.ID)
		return req, nil
	})
}

func (p *Proxy) HandleRequestAsUsual() goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, nil
	})
}
