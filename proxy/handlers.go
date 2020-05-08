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
	"github.com/vaulty/proxy/storage"
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
		vault, err := p.findVault(ctx, req)
		if err != nil {
			ctx.Warnf(err.Error())
			return nil, errResponse(req, "Vault was not found", http.StatusNotFound)
		}

		// for inbound requests req contains URL of host where
		// Vault is running e.g. proxy.vaulty.co. To forward
		// requests to vault's upstream we should set
		// host, port and user accordingly.
		if ctxUserData(ctx).routeType == model.RouteInbound {
			req.URL.Scheme = vault.UpstreamURL().Scheme
			req.URL.User = vault.UpstreamURL().User
			req.URL.Host = vault.UpstreamURL().Host
		}

		route, err := p.storage.FindRoute(vault.ID, ctxUserData(ctx).routeType, req)
		if err == storage.ErrNoRows {
			ctx.Warnf("Route was not found")
			return req, nil
		}
		if err != nil {
			ctx.Warnf(err.Error())
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		err = p.transformer.TransformRequestBody(route, req)
		if err != nil {
			return nil, errResponse(req, err.Error(), http.StatusInternalServerError)
		}

		ctxUserData(ctx).route = route

		return req, nil
	})
}

func (p *Proxy) findVault(ctx *goproxy.ProxyCtx, req *http.Request) (*model.Vault, error) {
	var (
		vaultID string
		err     error
	)

	// just return first vault if IsSingleVaultMode set
	if p.config.IsSingleVaultMode {
		vaults, err := p.storage.ListVaults()
		if err != nil {
			return nil, err
		}

		if len(vaults) == 0 {
			return nil, errors.New("No vaults found. Please, create vault first")
		}

		return vaults[0], nil
	}

	if ctxUserData(ctx).routeType == model.RouteInbound {
		vaultID, err = getVaultIDFromHost(p.config.BaseHost, req.Host)
		if err != nil {
			return nil, err
		}
	} else {
		vaultID = ctxUserData(ctx).vaultID
	}

	return p.storage.FindVault(vaultID)
}

var vaultIDRegexp *regexp.Regexp

// getVaultIDFromHost returns ID of vault from the host name
// We expect that in multi vault setup when base host is proxy.vaulty.co
// host name of vault looks like this:
// 	vltXXXX.proxy.vaulty.co
// 	vltZZZZ.proxy.vaulty.co
//
// Having proxy.vaulty.co as baseHost function returns vltXXXX
func getVaultIDFromHost(baseHost, host string) (string, error) {
	if vaultIDRegexp == nil {
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
		if ctxUserData(ctx).route == nil {
			return res
		}

		err := p.transformer.TransformResponseBody(ctxUserData(ctx).route, res)
		if err != nil {
			return errResponse(res.Request, err.Error(), http.StatusInternalServerError)
		}

		return res
	})
}

func (p *Proxy) HandleConnect() goproxy.HttpsHandler {
	return goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		vaultID, password, ok := proxyAuth(ctx.Req)

		if !ok || password != p.config.ProxyPassword {
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
