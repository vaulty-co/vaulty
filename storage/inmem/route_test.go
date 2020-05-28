package inmem

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/model"
	"github.com/vaulty/vaulty/storage"
)

func TestWithRoute(t *testing.T) {
	rs := NewStorage()
	defer rs.Reset()

	createdRoute := &model.Route{
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  "vlt1",
		Upstream: "http://example.com",
	}
	err := rs.CreateRoute(createdRoute)
	require.NoError(t, err)

	outboundRoute := &model.Route{
		Type:     model.RouteOutbound,
		Method:   http.MethodPost,
		Path:     "https://api.stripe.com/tokenize",
		VaultID:  "vlt1",
		Upstream: "http://example.com",
	}
	err = rs.CreateRoute(outboundRoute)
	require.NoError(t, err)

	t.Run("FindRoute", func(t *testing.T) {
		// inbound route
		req, err := http.NewRequest(http.MethodPost, "/tokenize", nil)
		require.NoError(t, err)

		route, err := rs.FindRoute("vlt1", model.RouteInbound, req)
		require.NoError(t, err)
		require.NotEmpty(t, route.ID)
		require.Equal(t, "http://example.com", route.Upstream)

		req, err = http.NewRequest(http.MethodPost, "/nothing", nil)
		require.NoError(t, err)

		route, err = rs.FindRoute("vlt1", model.RouteInbound, req)
		require.Equal(t, storage.ErrNoRows, err)

		// outbound route
		req, err = http.NewRequest(http.MethodPost, "https://api.stripe.com/tokenize", nil)
		require.NoError(t, err)

		route, err = rs.FindRoute("vlt1", model.RouteOutbound, req)
		require.NoError(t, err)
		require.NotEmpty(t, route.ID)
		require.Equal(t, outboundRoute, route)

	})

	t.Run("FindRouteByID", func(t *testing.T) {
		got, err := rs.FindRouteByID(createdRoute.VaultID, createdRoute.ID)
		require.NoError(t, err)
		require.Equal(t, createdRoute, got)

		_, err = rs.FindRouteByID("vlt1", "nothing")
		require.Equal(t, storage.ErrNoRows, err)
	})

	t.Run("ListRoutes", func(t *testing.T) {
		routes, err := rs.ListRoutes("vlt1")
		require.NoError(t, err)
		require.Len(t, routes, 2)
	})

	t.Run("DeleteRoute", func(t *testing.T) {
		err = rs.DeleteRoute(createdRoute.VaultID, createdRoute.ID)
		require.NoError(t, err)

		_, err = rs.FindRouteByID(createdRoute.VaultID, createdRoute.ID)
		require.Equal(t, storage.ErrNoRows, err)

		routes, err := rs.ListRoutes(createdRoute.VaultID)
		require.NoError(t, err)
		require.Len(t, routes, 1)
	})

	t.Run("DeleteRoutes", func(t *testing.T) {
		err := rs.CreateRoute(&model.Route{
			Type:     model.RouteInbound,
			Method:   http.MethodPost,
			Path:     "/tokenize1",
			VaultID:  "vlt1",
			Upstream: "http://example.com",
		})
		require.NoError(t, err)

		routes, err := rs.ListRoutes("vlt1")
		require.NoError(t, err)
		require.Len(t, routes, 2)

		err = rs.DeleteRoutes("vlt1")
		require.NoError(t, err)

		routes, err = rs.ListRoutes("vlt1")
		require.NoError(t, err)
		require.Len(t, routes, 0)
	})
}
