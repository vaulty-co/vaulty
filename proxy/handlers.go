package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/model"
)

func (p *Proxy) HandleRequest() goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		// find vault
		vaultID, err := getVaultID(p.config.BaseHost, req.Host)
		if err != nil {
			ctx.Warnf(err.Error())
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		vault, err := p.storage.FindVault(vaultID)
		if err != nil {
			ctx.Warnf(err.Error())
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		if vault == nil {
			return nil, errResponse(req, "Vault was not found", http.StatusNotFound)
		}

		req.URL.Scheme = vault.UpstreamURL().Scheme
		req.URL.User = vault.UpstreamURL().User
		req.URL.Host = vault.UpstreamURL().Host

		ctxUserData(ctx).vault = vault

		// find route
		route, err := p.storage.FindRoute(vault.ID, model.RouteInbound, req.Method, req.URL.Path)
		if err != nil {
			ctx.Warnf(err.Error())
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		if route == nil {
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
