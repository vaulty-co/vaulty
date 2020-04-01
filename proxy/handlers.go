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

func (p *Proxy) HandleRequest(message string) goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		// check if vault exists
		// check if route exists
		// transform request body
		// done
		return req, nil
	})
}

func (p *Proxy) HandleResponse() goproxy.RespHandler {
	return goproxy.FuncRespHandler(func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		// no route for current request
		if ctxUserData(ctx).route == nil {
			return res
		}

		err := p.transformer.TransformResponseBody(ctxUserData(ctx).route.ID, res)
		if err != nil {
			return errResponse(res.Request, err.Error(), http.StatusInternalServerError)
		}

		return res
	})
}
