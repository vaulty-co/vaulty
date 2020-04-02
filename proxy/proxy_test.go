package proxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage/test_storage"
	"github.com/vaulty/proxy/transformer/test_transformer"
)

type EchoHandler struct{}

func (EchoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, readBody(req.Body)+" response")
}

var upstream = httptest.NewServer(EchoHandler{})

func TestInboundRoute(t *testing.T) {
	defer test_storage.Reset()

	st := test_storage.NewTestStorage()
	tr := test_transformer.NewTransformer()
	config := core.LoadConfig("../config/test.yml")

	proxy := httptest.NewServer(NewProxy(st, tr, config).server)
	defer proxy.Close()

	test_storage.AddTestVault("vlt1", upstream.URL)

	err := st.CreateRoute(&model.Route{
		ID:       "rt1",
		Type:     model.RouteInbound,
		Method:   http.MethodPost,
		Path:     "/tokenize",
		VaultID:  "vlt1",
		Upstream: upstream.URL,
	})
	assert.NoError(t, err)

	t.Run("Test request and response body transformation when route matches", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, proxy.URL+"/tokenize", bytes.NewBufferString("request"))
		req.Host = "vlt1.proxy.test"

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
		req, _ := http.NewRequest(http.MethodPost, proxy.URL+"/pass", bytes.NewBufferString("request"))
		req.Host = "vlt1.proxy.test"

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

func readBody(body io.ReadCloser) string {
	b, err := ioutil.ReadAll(body)
	if err == nil {
		return string(b)
	}

	log.Fatal(err)

	return ""
}
