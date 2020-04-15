package storage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/model"
)

func TestFindRoute(t *testing.T) {
	rs := NewRedisStorage(redisClient)
	defer redisClient.FlushAll()

	err := rs.CreateRoute(&model.Route{
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  "vlt1",
		Upstream: "http://example.com",
	})
	require.NoError(t, err)

	route, err := rs.FindRoute("vlt1", model.RouteInbound, http.MethodPost, "/tokenize")
	require.NoError(t, err)
	require.NotNil(t, route)

	require.NotEmpty(t, route.ID)
	require.Equal(t, "http://example.com", route.Upstream)

	route, err = rs.FindRoute("vlt1", model.RouteInbound, http.MethodPost, "/nothing")
	require.Error(t, ErrNoRows, err)
}

func TestFindRouteByID(t *testing.T) {
	rs := NewRedisStorage(redisClient)
	defer redisClient.FlushAll()

	route := &model.Route{
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  "vlt1",
		Upstream: "http://example.com",
	}
	err := rs.CreateRoute(route)
	require.NoError(t, err)

	got, err := rs.FindRouteByID(route.VaultID, route.ID)
	require.NoError(t, err)
	require.Equal(t, route.Type, got.Type)
	require.Equal(t, route.Method, got.Method)
	require.Equal(t, route.Path, got.Path)
	require.Equal(t, route.Upstream, got.Upstream)
	require.Equal(t, route.RequestTransformationsJSON, got.RequestTransformationsJSON)
	require.Equal(t, route.ResponseTransformationsJSON, got.ResponseTransformationsJSON)

	route, err = rs.FindRouteByID("vlt1", "nothing")
	require.Equal(t, ErrNoRows, err)
}

func TestDeleteRoute(t *testing.T) {
	rs := NewRedisStorage(redisClient)
	defer redisClient.FlushAll()

	route := &model.Route{
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  "vlt1",
		Upstream: "http://example.com",
	}
	err := rs.CreateRoute(route)
	require.NoError(t, err)

	err = rs.DeleteRoute(route.VaultID, route.ID)
	require.NoError(t, err)

	_, err = rs.FindRouteByID(route.VaultID, route.ID)
	require.Equal(t, ErrNoRows, err)

	routes, err := rs.ListRoutes(route.VaultID)
	require.Len(t, routes, 0)
}
