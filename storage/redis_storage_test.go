package storage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vaulty/proxy/model"
)

func TestFindRoute(t *testing.T) {
	rs := NewRedisStorage(redisClient)

	err := rs.CreateRoute(&model.Route{
		ID:       "rt1",
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  "vlt1",
		Upstream: "http://example.com",
	})
	assert.NoError(t, err)

	t.Run("Finds route", func(t *testing.T) {
		route, err := rs.FindRoute("vlt1", model.RouteInbound, http.MethodPost, "/tokenize")

		assert.NoError(t, err)
		assert.Equal(t, route.ID, "rt1")
	})
}
