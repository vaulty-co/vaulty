package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/transformer"
)

func errResponse(r *http.Request, message string) *http.Response {
	return goproxy.NewResponse(r,
		goproxy.ContentTypeText,
		http.StatusBadGateway,
		message)
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		tr := transformer.NewRequestBodyTransformer(req)
		err := tr.TransformRequestBody()
		if err != nil {
			return nil, errResponse(req, err.Error())
		}

		return req, nil
	})

	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
