package proxy

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage/test_storage"
	"github.com/vaulty/proxy/transformer/test_transformer"
)

type EchoHandler struct{}

func (EchoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, readBody(req.Body)+" response")
}

var upstream = httptest.NewTLSServer(EchoHandler{})

func TestInboundRoute(t *testing.T) {
	defer test_storage.Reset()

	st := test_storage.NewTestStorage()
	tr := test_transformer.NewTransformer()
	config := core.LoadConfig("../config/test.yml")

	ps, err := NewProxy(st, tr, config)
	require.NoError(t, err)

	proxy := httptest.NewServer(ps.server)
	defer proxy.Close()

	vault := &model.Vault{
		ID:       "vlt1",
		Upstream: upstream.URL,
	}
	err = st.CreateVault(vault)
	require.NoError(t, err)

	err = st.CreateRoute(&model.Route{
		ID:       "rt1",
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  vault.ID,
		Upstream: upstream.URL,
	})
	require.NoError(t, err)

	t.Run("Test request and response body transformation when route matches", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, proxy.URL+"/tokenize", bytes.NewBufferString("request"))
		req.Host = fmt.Sprintf("%s.proxy.test", vault.ID)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		want := "request transformed response transformed"
		got := readBody(res.Body)

		if got != want {
			t.Errorf("Expected: %v, but got: %v", want, got)
		}
	})

	t.Run("Test request passes through to the vault's upstream when no route matches", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, proxy.URL+"/noroute", bytes.NewBufferString("request"))
		req.Host = fmt.Sprintf("%s.proxy.test", vault.ID)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		want := "request response"
		got := readBody(res.Body)

		if got != want {
			t.Errorf("Expected: %v, but got: %v", want, got)
		}
	})

	t.Run("Test request is rejected when no vault found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, proxy.URL+"/pass", bytes.NewBufferString("request"))
		req.Host = "vltunknown.proxy.test"

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != 404 {
			t.Errorf("Expected: %v, but got: %v", 404, res.StatusCode)
		}

		want := "Vault was not found"
		got := readBody(res.Body)

		if got != want {
			t.Errorf("Expected: %v, but got: %v", want, got)
		}
	})
}

func TestOutboundRoute(t *testing.T) {
	defer test_storage.Reset()

	st := test_storage.NewTestStorage()
	tr := test_transformer.NewTransformer()
	config := core.LoadConfig("../config/test.yml")

	ps, err := NewProxy(st, tr, config)
	require.NoError(t, err)

	proxy := httptest.NewServer(ps.server)
	defer proxy.Close()

	// example.com will never be reached for Outbound routes
	vault := &model.Vault{
		ID:       "vlt1",
		Upstream: "https://example.com",
	}
	err = st.CreateVault(vault)
	require.NoError(t, err)

	err = st.CreateRoute(&model.Route{
		ID:      "rt1",
		Type:    model.RouteOutbound,
		Method:  http.MethodPost,
		Path:    upstream.URL + "/tokenize",
		VaultID: vault.ID,
	})
	require.NoError(t, err)

	caCert, err := ioutil.ReadFile(filepath.Join(config.CaPath, "ca.pem"))
	require.NoError(t, err)

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	require.True(t, ok)

	tlsConfig := &tls.Config{}
	tlsConfig.RootCAs = caCertPool

	t.Run("Test proxy requires vault ID and pass in BasicAuth", func(t *testing.T) {
		// no user:password set for proxy
		transport := &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(proxy.URL)
			},
			TLSClientConfig: tlsConfig,
		}

		client := &http.Client{
			Transport: transport,
		}

		req, _ := http.NewRequest(http.MethodPost, upstream.URL+"/tokenize", bytes.NewBufferString("request"))

		_, err := client.Do(req)
		require.Contains(t, err.Error(), "Proxy Authentication Required")
	})

	t.Run("Test request and response body transformation when route matches", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, upstream.URL+"/tokenize", bytes.NewBufferString("request"))

		proxyURL, _ := url.Parse(proxy.URL)
		proxyURL.User = url.UserPassword(vault.ID, config.ProxyPassword)

		transport := &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return proxyURL, nil
			},
			TLSClientConfig: tlsConfig,
		}

		client := &http.Client{
			Transport: transport,
		}

		res, err := client.Do(req)
		require.NoError(t, err)

		want := "request transformed response transformed"
		got := readBody(res.Body)

		if got != want {
			t.Errorf("Expected: %v, but got: %v", want, got)
		}
	})
}

func readBody(body io.ReadCloser) string {
	b, err := ioutil.ReadAll(body)
	if err == nil {
		return string(b)
	}

	log.Fatal(err)

	return ""
}
