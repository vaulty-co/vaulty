package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/model"
)

func (p *Proxy) vaultDoesNotExist() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		vaultID, err := getVaultID(req.Host)
		if err != nil {
			ctx.Warnf(err.Error())
			return true
		}

		vault, err := p.storage.FindVault(vaultID)
		if err != nil {
			ctx.Warnf(err.Error())
			return true
		}

		if vault == nil {
			return true
		}

		req.URL = vault.UpstreamURL

		ctxUserData(ctx).vault = vault

		return false
	}
}

var vaultIDRegexp *regexp.Regexp

func getVaultID(host string) (string, error) {
	if vaultIDRegexp == nil {
		// vltXXXX.proxy.vaulty.co
		vaultHost := fmt.Sprintf(`^(vlt\w+).%s(:\d+)?$`, core.Config.BaseHost)
		vaultIDRegexp = regexp.MustCompile(vaultHost)
	}

	matches := vaultIDRegexp.FindAllStringSubmatch(host, -1)

	if len(matches) != 1 {
		return "", errors.New(fmt.Sprintf("Received request for %s instead of configured host: vlt*.%s", host, core.Config.BaseHost))
	}

	vaultID := matches[0][1]

	return vaultID, nil
}

// matches route and find route id
// vlt2uYBrnYkUnEF:INBOUND:POST:/records
func (p *Proxy) routeExists() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		vaultID := ctxUserData(ctx).vault.ID
		route, err := p.storage.FindRoute(vaultID, model.RouteInbound, req.Method, req.URL.Path)
		if err != nil {
			ctx.Warnf(err.Error())
			return false
		}

		if route == nil {
			return false
		}

		ctxUserData(ctx).route = route

		return true
	}
}
