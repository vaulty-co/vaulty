package proxy

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/storage"
)

type Proxy struct {
	server  *goproxy.ProxyHttpServer
	storage *storage.Storage
}

func NewProxy() *Proxy {
	server := goproxy.NewProxyHttpServer()
	redisClient := core.NewRedisClient()
	storage := storage.NewStorage(redisClient)

	proxy := &Proxy{
		server:  server,
		storage: storage,
	}

	server.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = "inbound-request.int"

		server.ServeHTTP(w, req)
	})

	// proxy.OnRequest(matchOutboundRoute()).HandleConnect(goproxy.AlwaysMitm)

	server.OnRequest(proxy.vaultDoesNotExist()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return nil, errResponse(req, "Vault was not found", http.StatusNotFound)
	})

	server.OnRequest(proxy.routeExists()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		fmt.Println("route found: ", ctxUserData(ctx).route.ID)
		return req, nil
	})

	// proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 	fmt.Println("Global handler")
	// 	return req, nil
	// })

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
