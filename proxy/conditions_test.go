package proxy

import (
	"net/http/httptest"
	"testing"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/model"
)

func TestVaultDoesNotExist(t *testing.T) {
	proxy := NewProxy()
	condition := proxy.vaultDoesNotExist()

	redisClient.Set("vault:vlt123:upstream", "https://upstream.com", 0)

	t.Run("Vault Exists", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://vlt123.proxy.test/foo", nil)
		ctx := &goproxy.ProxyCtx{}

		res := condition(req, ctx)

		if res == true {
			t.Error("Vault was not found")
		}

		// if ctxuserdata(ctx).vault.upstreamurl.
		if ctxUserData(ctx).vault.UpstreamURL.Host != "upstream.com" {
			t.Error("Upstream was not set")
		}
	})

	t.Run("Vault Does Not Exist", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://vlt456.proxy.test/foo", nil)
		ctx := &goproxy.ProxyCtx{}

		res := condition(req, ctx)

		if res == false {
			t.Error("Vault was found")
		}
	})

}

func TestRouteExist(t *testing.T) {
	proxy := NewProxy()
	condition := proxy.routeExists()

	// create route with upstream
	redisClient.Set("vlt123:INBOUND:POST:/foo", "1", 0)
	redisClient.Set("route:1:upstream", "https://upstream.com", 0)

	t.Run("Route Exists", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://vlt123.proxy.test/foo", nil)
		ctx := &goproxy.ProxyCtx{}
		ctxUserData(ctx).vault = &model.Vault{
			ID: "vlt123",
		}

		res := condition(req, ctx)

		if res == false {
			t.Error("Route was not found")
		}

		route := ctxUserData(ctx).route

		if (route.ID != "1") || (route.UpstreamURL.Host != "upstream.com") {
			t.Error("Route with ID and UpstreamURL was not found")
		}
	})

	t.Run("Route Does Not Exist", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://vlt456.proxy.test/foo", nil)
		ctx := &goproxy.ProxyCtx{}
		ctxUserData(ctx).vault = &model.Vault{
			ID: "vlt123",
		}

		res := condition(req, ctx)

		if res == false {
			t.Error("Route was found")
		}
	})
}
