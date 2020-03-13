package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
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

		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, errResponse(req, err.Error())
		}

		oldBody := string(buf)
		newBody := " + new body is here!!!"

		buf2 := bytes.NewBufferString(oldBody + newBody)
		req.Header.Del("Content-Length")
		req.ContentLength = int64(buf2.Len())
		req.Body = ioutil.NopCloser(bufio.NewReader(buf2))

		return req, nil
	})

	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
