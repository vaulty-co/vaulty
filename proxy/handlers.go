package proxy

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"github.com/vaulty/proxy/model"
)

func (p *Proxy) SetRouteType() goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if ctxUserData(ctx).routeType == "" {
			ctxUserData(ctx).routeType = model.RouteInbound
		}

		return req, nil
	})
}
func (p *Proxy) HandleRequest() goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		var (
			vaultID string
			err     error
		)

		if ctxUserData(ctx).routeType == model.RouteInbound {
			vaultID, err = getVaultID(p.config.BaseHost, req.Host)
			if err != nil {
				ctx.Warnf(err.Error())
				return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
			}
		} else {
			vaultID = ctxUserData(ctx).vaultID
		}

		vault, err := p.storage.FindVault(vaultID)
		if err != nil {
			ctx.Warnf(err.Error())
			return nil, errResponse(req, "Vault was not found", http.StatusNotFound)
		}

		ctxUserData(ctx).vault = vault

		if ctxUserData(ctx).routeType == model.RouteInbound {
			req.URL.Scheme = vault.UpstreamURL().Scheme
			req.URL.User = vault.UpstreamURL().User
			req.URL.Host = vault.UpstreamURL().Host
		}

		// find route
		route, err := p.storage.FindRoute(vault.ID, ctxUserData(ctx).routeType, req.Method, req.URL.Path)
		if err != nil {
			ctx.Warnf(err.Error())
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		if route == nil {
			ctx.Warnf("No route found")
			return req, nil
		}

		ctxUserData(ctx).route = route

		// transform request
		err = p.transformer.TransformRequestBody(ctxUserData(ctx).route, req)
		if err != nil {
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		return req, nil
	})
}

var vaultIDRegexp *regexp.Regexp

func getVaultID(baseHost, host string) (string, error) {
	if vaultIDRegexp == nil {
		// vltXXXX.proxy.vaulty.co
		vaultHost := fmt.Sprintf(`^(vlt\w+).%s(:\d+)?$`, baseHost)
		vaultIDRegexp = regexp.MustCompile(vaultHost)
	}

	matches := vaultIDRegexp.FindAllStringSubmatch(host, -1)

	if len(matches) != 1 {
		return "", errors.New(fmt.Sprintf("Received request for %s instead of configured host: vlt*.%s", host, baseHost))
	}

	vaultID := matches[0][1]

	return vaultID, nil
}

func (p *Proxy) HandleResponse() goproxy.RespHandler {
	return goproxy.FuncRespHandler(func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		// no route for current request
		if ctxUserData(ctx).route == nil {
			return res
		}

		// transform response
		err := p.transformer.TransformResponseBody(ctxUserData(ctx).route, res)
		if err != nil {
			return errResponse(res.Request, err.Error(), http.StatusInternalServerError)
		}

		return res
	})
}

func (p *Proxy) HandleConnect() goproxy.HttpsHandler {
	return goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		// handle error here
		// if no vaultID
		vaultID, password, ok := proxyAuth(ctx.Req)

		if !ok || password != "pass" {
			ctx.Resp = auth.BasicUnauthorized(ctx.Req, "")
			return goproxy.RejectConnect, host
		}

		ctxUserData(ctx).routeType = model.RouteOutbound
		ctxUserData(ctx).vaultID = vaultID

		return goproxy.MitmConnect, host
	})
}

var proxyAuthorizationHeader = "Proxy-Authorization"

func proxyAuth(req *http.Request) (user, passwd string, ok bool) {
	authheader := strings.SplitN(req.Header.Get("Proxy-Authorization"), " ", 2)
	// req.Header.Del("Proxy-Authorization")
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
