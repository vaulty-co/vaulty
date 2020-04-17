package storage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/model"
)

func TestWithRoute(t *testing.T) {
	rs := NewRedisStorage(redisClient)
	defer redisClient.FlushAll()

	createdRoute := &model.Route{
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  "vlt1",
		Upstream: "http://example.com",
	}
	err := rs.CreateRoute(createdRoute)
	require.NoError(t, err)

	t.Run("FindRoute", func(t *testing.T) {
		route, err := rs.FindRoute("vlt1", model.RouteInbound, http.MethodPost, "/tokenize")
		require.NoError(t, err)
		require.NotEmpty(t, route.ID)
		require.Equal(t, "http://example.com", route.Upstream)

		route, err = rs.FindRoute("vlt1", model.RouteInbound, http.MethodPost, "/nothing")
		require.Equal(t, ErrNoRows, err)
	})

	t.Run("FindRouteByID", func(t *testing.T) {
		got, err := rs.FindRouteByID(createdRoute.VaultID, createdRoute.ID)
		require.NoError(t, err)
		require.Equal(t, createdRoute, got)

		_, err = rs.FindRouteByID("vlt1", "nothing")
		require.Equal(t, ErrNoRows, err)
	})

	t.Run("DeleteRoute", func(t *testing.T) {
		err = rs.DeleteRoute(createdRoute.VaultID, createdRoute.ID)
		require.NoError(t, err)

		_, err = rs.FindRouteByID(createdRoute.VaultID, createdRoute.ID)
		require.Equal(t, ErrNoRows, err)

		routes, err := rs.ListRoutes(createdRoute.VaultID)
		require.NoError(t, err)
		require.Len(t, routes, 0)
	})
}
