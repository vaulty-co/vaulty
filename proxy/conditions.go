package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/redis"
)

func vaultDoesNotExist() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		fmt.Println(req.Header)

		vaultID, err := getVaultID(req.Host)
		if err != nil {
			ctx.Warnf(err.Error())
			return true
		}

		vlt := &vault{
			ID: vaultID,
		}

		upstreamURL, err := url.Parse(vlt.GetUpstream())

		fmt.Println("Got upstreamURL: ", upstreamURL)
		if err != nil {
			ctx.Warnf(err.Error())
			return true
		}

		req.URL = upstreamURL

		ctxUserData(ctx).vault = vlt

		return false
	}
}

var vaultIDRegexp *regexp.Regexp

func getVaultID(host string) (string, error) {
	if vaultIDRegexp == nil {
		// vltXXXX.proxy.vaulty.co
		vaultHost := fmt.Sprintf(`^(vlt\w+).%s(:\d+)?$`, core.Config().Host)
		vaultIDRegexp = regexp.MustCompile(vaultHost)
	}

	matches := vaultIDRegexp.FindAllStringSubmatch(host, -1)

	if len(matches) != 1 {
		return "", errors.New(fmt.Sprintf("Received request for %s instead of configured host: vlt*.%s", host, core.Config().Host))
	}

	vaultID := matches[0][1]

	return vaultID, nil
}

// matches route and find route id
// vlt2uYBrnYkUnEF:INBOUND:POST:/records
func routeExists() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		vaultID := ctxUserData(ctx).vault.ID
		routeKey := fmt.Sprintf("%s:%s:%s:%s", vaultID, "INBOUND", req.Method, req.URL.Path)
		ctx.Logf("Route key: " + routeKey)

		routeID := redis.Client().Get(routeKey).Val()
		ctx.Logf("RouteID: " + routeID)

		if routeID == "" {
			ctx.Logf("Route was not found")
			return false
		}

		ctxUserData(ctx).route = &route{
			ID: routeID,
		}

		return true
	}
}
