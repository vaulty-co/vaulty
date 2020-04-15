package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/api/request"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage/test_storage"
)

func TestHandleRouteCreate(t *testing.T) {
	st := test_storage.NewTestStorage()
	server := NewServer(test_storage.NewTestStorage())
	defer test_storage.Reset()

	vault := &model.Vault{Upstream: "https://example.com"}
	err := st.CreateVault(vault)
	require.Nil(t, err)

	in := new(bytes.Buffer)
	json.NewEncoder(in).Encode(&model.Route{
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		Upstream: "https://example.com",
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)
	r = r.WithContext(request.WithVault(r.Context(), vault))

	server.HandleRouteCreate()(w, r)

	require.Equal(t, 200, w.Code)

	out := &model.Route{}
	json.NewDecoder(w.Body).Decode(out)

	require.NotEmpty(t, out.ID)
	require.Equal(t, model.RouteInbound, out.Type)
	require.Equal(t, "POST", out.Method)
	require.Equal(t, "/tokenize", out.Path)
	require.Equal(t, "https://example.com", out.Upstream)
	require.Equal(t, vault.ID, out.VaultID)
}
