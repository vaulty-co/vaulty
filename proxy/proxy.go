package proxy

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

type Proxy struct {
	server *goproxy.ProxyHttpServer
}

func NewProxy() *Proxy {
	proxy := goproxy.NewProxyHttpServer()

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = "inbound-request.int"

		proxy.ServeHTTP(w, req)
	})

	// proxy.OnRequest(matchOutboundRoute()).HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest(vaultDoesNotExist()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return nil, errResponse(req, "Vault was not found", http.StatusNotFound)
	})

	proxy.OnRequest(routeExists()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		fmt.Println("route found: ", ctxUserData(ctx).route.ID)
		return req, nil
	})

	// proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 	fmt.Println("Global handler")
	// 	return req, nil
	// })

	proxy.Verbose = true

	return &Proxy{
		server: proxy,
	}
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
