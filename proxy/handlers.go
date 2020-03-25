package proxy

import (
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
		ctx.Logf("Round found, transform data")

		err := p.transformer.TransformRequestBody(ctxUserData(ctx).route.ID, req)
		if err != nil {
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		return req, nil
	})
}

func (p *Proxy) HandleRequestAsUsual() goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, nil
	})
}
