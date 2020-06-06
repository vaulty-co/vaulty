package proxy

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	log "github.com/sirupsen/logrus"
	"github.com/vaulty/vaulty/routing"
)

func (p *Proxy) HandleRequest() goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

		route := p.router.LookupRoute(req)

		if route == nil {
			return nil, errResponse(req, "No route found", http.StatusNotFound)
		}

		log.Debugf("Route found: %s", route.Name)

		if route.IsInbound {
			req.URL.Scheme = route.UpstreamURL.Scheme
			req.URL.User = route.UpstreamURL.User
			req.URL.Host = route.UpstreamURL.Host
		}

		req, err := p.TransformRequestBody(route, req)
		if err != nil {
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		ctxUserData(ctx).route = route

		return req, nil
	})
}

func (p *Proxy) HandleResponse() goproxy.RespHandler {
	return goproxy.FuncRespHandler(func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if ctxUserData(ctx).route == nil {
			return res
		}

		if res == nil {
			return res
		}

		res, err := p.TransformResponseBody(ctxUserData(ctx).route, res)
		if err != nil {
			return errResponse(res.Request, err.Error(), http.StatusInternalServerError)
		}

		return res
	})
}

func (p *Proxy) TransformRequestBody(route *routing.Route, req *http.Request) (*http.Request, error) {
	req, err := route.TransformRequest(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (p *Proxy) TransformResponseBody(route *routing.Route, res *http.Response) (*http.Response, error) {
	res, err := route.TransformResponse(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Proxy) HandleConnect() goproxy.HttpsHandler {
	return goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		_, password, ok := proxyAuth(ctx.Req)

		if !ok || password != p.proxyPassword {
			ctx.Resp = auth.BasicUnauthorized(ctx.Req, "")
			return goproxy.RejectConnect, host
		}

		return goproxy.MitmConnect, host
	})
}

var proxyAuthorizationHeader = "Proxy-Authorization"

func proxyAuth(req *http.Request) (user, passwd string, ok bool) {
	authheader := strings.SplitN(req.Header.Get("Proxy-Authorization"), " ", 2)
	if len(authheader) != 2 || authheader[0] != "Basic" {
		return "", "", false
	}
	userpassraw, err := base64.StdEncoding.DecodeString(authheader[1])
	if err != nil {
		return "", "", false
	}
	userpass := strings.SplitN(string(userpassraw), ":", 2)
	if len(userpass) != 2 {
		return "", "", false
	}
	return userpass[0], userpass[1], true
}
