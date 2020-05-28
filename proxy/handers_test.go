package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elazarl/goproxy"
	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/model"
	"github.com/vaulty/vaulty/storage/inmem"
)

func TestHandleRequest(t *testing.T) {
	st := inmem.NewStorage()
	defer st.Reset()

	opts := &Options{
		CAPath:  "./testdata",
		Storage: st,
	}

	proxy, err := NewProxy(opts)
	require.NoError(t, err)

	vault := &model.Vault{
		Upstream: "https://default-backend.com",
	}
	err = st.CreateVault(vault)
	require.NoError(t, err)

	ctx := &goproxy.ProxyCtx{}
	ctxUserData(ctx).routeType = model.RouteInbound

	t.Run("Test default upstream is used when route's upstream is not set", func(t *testing.T) {
		err = st.CreateRoute(&model.Route{
			Type:    model.RouteInbound,
			Method:  http.MethodPost,
			Path:    "/without-upstream",
			VaultID: vault.ID,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "https://proxy.vaulty.co/without-upstream", nil)
		handler := proxy.HandleRequest()

		req, _ = handler.Handle(req, ctx)

		require.Equal(t, "https://default-backend.com/without-upstream", req.URL.String())
	})

	t.Run("Test route's upstream is used when set", func(t *testing.T) {
		err = st.CreateRoute(&model.Route{
			Type:     model.RouteInbound,
			Method:   http.MethodPost,
			Path:     "/with-upstream",
			VaultID:  vault.ID,
			Upstream: "https://route-backend.com",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "https://proxy.vaulty.co/with-upstream", nil)
		handler := proxy.HandleRequest()

		req, _ = handler.Handle(req, ctx)

		require.Equal(t, "https://route-backend.com/with-upstream", req.URL.String())
	})
}
