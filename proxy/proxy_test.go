package proxy

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/routing"
	"github.com/vaulty/vaulty/transformer"
)

type EchoHandler struct{}

func (EchoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, readBody(req.Body)+" response")
}

var upstream = httptest.NewTLSServer(EchoHandler{})

type fakeTransformer struct{}

func (f *fakeTransformer) TransformRequest(req *http.Request) (*http.Request, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	body = append(body, " transformed"...)

	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	req.Header.Del("Content-Length")
	req.ContentLength = int64(len(body))

	return req, nil
}

func (f *fakeTransformer) TransformResponse(res *http.Response) (*http.Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	body = append(body, " transformed"...)

	res.Body = ioutil.NopCloser(bytes.NewReader(body))
	res.Header.Del("Content-Length")
	res.ContentLength = int64(len(body))

	return res, nil
}

func TestInboundRoute(t *testing.T) {
	route, err := routing.NewRoute(&routing.RouteParams{
		Name:                    "in",
		Method:                  http.MethodPost,
		URL:                     "/tokenize",
		Upstream:                upstream.URL,
		RequestTransformations:  []transformer.Transformer{&fakeTransformer{}},
		ResponseTransformations: []transformer.Transformer{&fakeTransformer{}},
	})
	require.NoError(t, err)

	router := routing.NewRouter()
	router.SetRoutes([]*routing.Route{route})

	opts := &Options{
		CAPath: "./testdata",
		Router: router,
	}

	ps, err := NewProxy(opts)
	require.NoError(t, err)

	proxy := httptest.NewServer(ps.server)
	defer proxy.Close()

	t.Run("Test request and response body transformation when route matches", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, proxy.URL+"/tokenize", bytes.NewBufferString("request"))

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		want := "request transformed response transformed"
		got := readBody(res.Body)

		require.Equal(t, want, got)
	})

	t.Run("Test proxy returns 404 if no route matches", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, proxy.URL+"/noroute", bytes.NewBufferString("request"))

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		want := "No route found"
		got := readBody(res.Body)

		require.Equal(t, want, got)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestOutboundRoute(t *testing.T) {
	route, err := routing.NewRoute(&routing.RouteParams{
		Name:                    "out",
		Method:                  http.MethodPost,
		URL:                     upstream.URL + "/tokenize",
		RequestTransformations:  []transformer.Transformer{&fakeTransformer{}},
		ResponseTransformations: []transformer.Transformer{&fakeTransformer{}},
	})
	require.NoError(t, err)

	router := routing.NewRouter()
	router.SetRoutes([]*routing.Route{route})

	opts := &Options{
		CAPath:        "./testdata",
		Router:        router,
		ProxyPassword: "password",
	}

	ps, err := NewProxy(opts)
	require.NoError(t, err)

	proxy := httptest.NewServer(ps.server)
	defer proxy.Close()

	caCert, err := ioutil.ReadFile(filepath.Join(opts.CAPath, "ca.cert"))
	require.NoError(t, err)

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	require.True(t, ok)

	tlsConfig := &tls.Config{}
	tlsConfig.RootCAs = caCertPool

	t.Run("Test proxy requires password in BasicAuth", func(t *testing.T) {
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

		// setup user:password for proxyURL to pass basic auth
		proxyURL, _ := url.Parse(proxy.URL)
		proxyURL.User = url.UserPassword("x", opts.ProxyPassword)

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

	t.Run("Test request and response body transformation when route matches", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, upstream.URL+"/tokenize", bytes.NewBufferString("request"))

		// basic auth password "x" should not be treated as vault ID
		// in signle vault mode
		proxyURL, _ := url.Parse(proxy.URL)
		proxyURL.User = url.UserPassword("x", opts.ProxyPassword)

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

	t.Run("Test proxy returns 404 if no route matches", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "https://unknown.com/path", bytes.NewBufferString("request"))

		proxyURL, _ := url.Parse(proxy.URL)
		proxyURL.User = url.UserPassword("x", opts.ProxyPassword)

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
		if err != nil {
			t.Error(err)
		}

		want := "No route found"
		got := readBody(res.Body)

		require.Equal(t, want, got)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
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
