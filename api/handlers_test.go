package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage/test_storage"
)

func TestVaultCreate(t *testing.T) {
	server := NewServer(test_storage.NewTestStorage())
	in := new(bytes.Buffer)
	json.NewEncoder(in).Encode(&model.Vault{Upstream: "https://example.com"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	server.HandleVaultCreate()(w, r)

	require.Equal(t, 200, w.Code)

	out := &model.Vault{}
	json.NewDecoder(w.Body).Decode(out)
	require.NotEmpty(t, out.ID)
	require.Equal(t, "https://example.com", out.Upstream)
}
