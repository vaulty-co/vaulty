package proxy

import (
	"log"
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

func NewProxy(storage storage.Storage, transformer transformer.Transformer, config *core.Configuration) *Proxy {
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

	// proxy.OnRequest(matchOutboundRoute()).HandleConnect(goproxy.AlwaysMitm)
	server.OnRequest().Do(proxy.HandleRequest())
	server.OnResponse().Do(proxy.HandleResponse())
	server.Verbose = true

	return proxy
}

func (p *Proxy) Run(port string) {
	log.Fatal(http.ListenAndServe(":"+port, p.server))
}
