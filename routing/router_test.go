package routing

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupRoute(t *testing.T) {
	in, err := NewRoute(&RouteParams{
		Method:   "POST",
		URL:      "/tokenize",
		Upstream: "http://backend",
	})
	require.NoError(t, err)

	out, err := NewRoute(&RouteParams{
		Method: "POST",
		URL:    "https://example.com/tokenize",
	})
	require.NoError(t, err)

	router := NewRouter()
	router.SetRoutes([]*Route{in, out})

	req := httptest.NewRequest("POST", "https://inbound/tokenize", nil)
	require.Equal(t, in, router.LookupRoute(req))

	req = httptest.NewRequest("POST", "https://example.com/tokenize", nil)
	require.Equal(t, out, router.LookupRoute(req))
}
