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
	storage     *storage.Storage
	transformer *transformer.Transformer
	config      *core.Configuration
}

func NewProxy(storage *storage.Storage, transformer *transformer.Transformer, config *core.Configuration) *Proxy {
	server := goproxy.NewProxyHttpServer()

	proxy := &Proxy{
		server:      server,
		storage:     storage,
		transformer: transformer,
		config:      config,
	}

	server.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = "inbound-request.int"

		server.ServeHTTP(w, req)
	})

	// proxy.OnRequest(matchOutboundRoute()).HandleConnect(goproxy.AlwaysMitm)

	// server.OnRequest().Do(proxy.FindVaultAndRoute())

	// if vault does not exist we respond with 404
	server.OnRequest(proxy.vaultDoesNotExist()).Do(proxy.NotFound("Vault was not found"))

	// if vault exist and there is a route for current request
	server.OnRequest(proxy.routeExists()).Do(proxy.TransformRequestBody())

	// if vault exist and there were no route
	server.OnRequest().Do(proxy.HandleRequestAsUsual())

	server.Verbose = true

	return proxy
}

func (p *Proxy) Run() {
	log.Fatal(http.ListenAndServe(":8080", p.server))
}

// match route and find route id
// vlt2uYBrnYkUnEF:INBOUND:POST:/records
// var routeID string

// if req.URL.Path == "/credit-cards" && req.Method == "POST" {
// 	routeID = "1"
// }

// tr := transformer.NewTransformer()
// err := tr.TransformRequestBody(routeID, req)
// if err != nil {
// 	return nil, errResponse(req, err.Error())
// }
