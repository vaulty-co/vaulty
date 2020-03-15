package main

import (
	"fmt"
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
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = "postman-echo.com"
		req.Host = "postman-echo.com"

		proxy.ServeHTTP(w, req)
	})
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		tr := transformer.NewTransformer()
		err := tr.TransformRequestBody(routeID, req)
		if err != nil {
			return nil, errResponse(req, err.Error())
		}

		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		fmt.Println("Response Status:", resp.StatusCode)
		fmt.Println("Response for request host:", resp.Request.Host)
		// resp.StatusCode = http.StatusOK
		// resp.Body = ioutil.NopCloser(bytes.NewBufferString("chico"))
		return resp
	})

	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
