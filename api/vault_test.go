package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/api/request"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
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

func TestVaultCtx(t *testing.T) {
	st := test_storage.NewTestStorage()
	server := NewServer(st)
	defer test_storage.Reset()

	vault := &model.Vault{Upstream: "https://example.com"}
	err := st.CreateVault(vault)
	require.Nil(t, err)

	t.Run("Test vault is set into the request context", func(t *testing.T) {
		routeCtx := new(chi.Context)
		routeCtx.URLParams.Add("vaultID", vault.ID)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r = r.WithContext(context.WithValue(request.WithVault(r.Context(), vault), chi.RouteCtxKey, routeCtx))

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, vault, request.VaultFrom(r.Context()))
		})

		server.VaultCtx(testHandler).ServeHTTP(w, r)
	})

	t.Run("Test vault not found", func(t *testing.T) {
		c := new(chi.Context)
		c.URLParams.Add("vaultID", "xxx")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, c))

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Should never be called")
		})

		server.VaultCtx(testHandler).ServeHTTP(w, r)

		require.Equal(t, 404, w.Code)
	})
}

func TestHandleVaultFind(t *testing.T) {
	st := test_storage.NewTestStorage()
	server := NewServer(st)
	defer test_storage.Reset()

	vault := &model.Vault{Upstream: "https://example.com"}
	err := st.CreateVault(vault)
	require.Nil(t, err)

	t.Run("Returns vault", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r = r.WithContext(request.WithVault(r.Context(), vault))

		server.HandleVaultFind()(w, r)

		require.Equal(t, 200, w.Code)

		out := &model.Vault{}
		json.NewDecoder(w.Body).Decode(out)
		require.Equal(t, vault.ID, out.ID)
	})
}

func TestHandleVaultUpdate(t *testing.T) {
	st := test_storage.NewTestStorage()
	server := NewServer(st)
	defer test_storage.Reset()

	vault := &model.Vault{Upstream: "https://example.com"}
	err := st.CreateVault(vault)
	require.Nil(t, err)

	t.Run("Updates vault", func(t *testing.T) {
		in := new(bytes.Buffer)
		json.NewEncoder(in).Encode(&vaultInput{Upstream: "https://newdomain.com"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/", in)
		r = r.WithContext(request.WithVault(r.Context(), vault))

		server.HandleVaultUpdate()(w, r)

		require.Equal(t, 200, w.Code)

		out := &model.Vault{}
		json.NewDecoder(w.Body).Decode(out)
		require.NotEmpty(t, out.ID)
		require.Equal(t, "https://newdomain.com", out.Upstream)
	})
}

func TestHandleVaultDelete(t *testing.T) {
	st := test_storage.NewTestStorage()
	server := NewServer(st)
	defer test_storage.Reset()

	vault := &model.Vault{Upstream: "https://example.com"}
	err := st.CreateVault(vault)
	require.Nil(t, err)

	// todo all routes
	t.Run("Deletes vault", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/", nil)
		r = r.WithContext(request.WithVault(r.Context(), vault))

		server.HandleVaultDelete()(w, r)

		require.Equal(t, 204, w.Code)

		_, err := st.FindVault(vault.ID)
		require.Equal(t, storage.ErrNoRows, err)
	})
}
