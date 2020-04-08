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

func TestHandleVaultCreate(t *testing.T) {
	server := NewServer(test_storage.NewTestStorage())
	defer test_storage.Reset()

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

func TestHandleVaultList(t *testing.T) {
	st := test_storage.NewTestStorage()
	server := NewServer(st)
	defer test_storage.Reset()

	vault := &model.Vault{Upstream: "https://example.com"}
	err := st.CreateVault(vault)
	require.Nil(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	server.HandleVaultList()(w, r)
	require.Equal(t, 200, w.Code)

	want := []*model.Vault{
		vault,
	}

	got := []*model.Vault{}
	json.NewDecoder(w.Body).Decode(&got)

	require.Equal(t, want, got)
}
