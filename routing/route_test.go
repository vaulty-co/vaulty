package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	t.Run("Test inbound route for specific path", func(t *testing.T) {
		route, err := NewRoute(&RouteParams{
			Method: "POST",
			URL:    "/tokenize",
		})
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "http://inbound/tokenize", nil)
		require.True(t, route.Match(req))
	})

	t.Run("Test inbound route for all requests", func(t *testing.T) {
		route, err := NewRoute(&RouteParams{
			Method: "*",
			URL:    "/*",
		})
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "http://inbound/whatever", nil)
		require.True(t, route.Match(req))
	})

	t.Run("Test outbound route for specific path", func(t *testing.T) {
		route, err := NewRoute(&RouteParams{
			Method: "POST",
			URL:    "https://api.stripe.com/tokenize",
		})
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "https://api.stripe.com/tokenize", nil)
		require.True(t, route.Match(req))
	})

	t.Run("Test outbound route for all requests", func(t *testing.T) {
		route, err := NewRoute(&RouteParams{
			Method: "*",
			URL:    "https://api.stripe.com/*",
		})
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "https://api.stripe.com/whatever", nil)
		require.True(t, route.Match(req))
	})

	t.Run("Test outbound route for request without ending /", func(t *testing.T) {
		route, err := NewRoute(&RouteParams{
			Method: "*",
			URL:    "https://api.stripe.com/*",
		})
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", "https://api.stripe.com", nil)
		require.True(t, route.Match(req))
	})
}

func TestIsInbound(t *testing.T) {
	var routeTests = []struct {
		url       string
		isInbound bool
	}{
		{"/tokenize", true},
		{"/*", true},
		{"http://example.com", false},
		{"http://example.com/", false},
		{"http://example.com/*", false},
	}

	for _, tt := range routeTests {
		t.Run(tt.url, func(t *testing.T) {
			route, err := NewRoute(&RouteParams{
				URL: tt.url,
			})

			require.NoError(t, err)
			require.Equal(t, tt.isInbound, route.IsInbound)
		})
	}
}
