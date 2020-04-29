package proxy

import (
	"errors"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/storage"
	"github.com/vaulty/proxy/transformer"
)

type Proxy struct {
	server      *goproxy.ProxyHttpServer
	storage     storage.Storage
	transformer transformer.Transformer
	config      *core.Configuration
}

func NewProxy(storage storage.Storage, transformer transformer.Transformer, config *core.Configuration) (*Proxy, error) {
	if config.ProxyPassword == "" {
		return nil, errors.New("Proxy password must be specified via config file or PROXY_PASS environment variable")
	}

	server := goproxy.NewProxyHttpServer()

	proxy := &Proxy{
		server:      server,
		storage:     storage,
		transformer: transformer,
		config:      config,
	}

	server.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = "inbound-request"

		server.ServeHTTP(w, req)
	})

	server.OnRequest().HandleConnect(proxy.HandleConnect())
	server.OnRequest().Do(proxy.SetRouteType())
	server.OnRequest().Do(proxy.HandleRequest())
	server.OnResponse().Do(proxy.HandleResponse())
	server.Verbose = true

	return proxy, nil
}

func (s *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}
