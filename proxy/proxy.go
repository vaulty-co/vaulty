package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/vaulty/routing"
)

type Options struct {
	// Password for the forward proxy
	ProxyPassword string

	// Path to CA files
	CAPath string

	// router with all routes
	Router routing.Router
}

type Proxy struct {
	proxyPassword string
	server        *goproxy.ProxyHttpServer
	router        routing.Router
}

func NewProxy(opts *Options) (*Proxy, error) {
	server := goproxy.NewProxyHttpServer()

	err := setupCA(opts.CAPath)
	if err != nil {
		return nil, err
	}

	proxy := &Proxy{
		server:        server,
		router:        opts.Router,
		proxyPassword: opts.ProxyPassword,
	}

	server.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// set host for inbound requests
		req.URL.Scheme = "https"
		req.URL.Host = "inbound"
		server.ServeHTTP(w, req)
	})

	server.OnRequest().HandleConnect(proxy.HandleConnect())
	server.OnRequest().Do(proxy.HandleRequest())
	server.OnResponse().Do(proxy.HandleResponse())
	server.Verbose = true

	return proxy, nil
}

func setupCA(CAPath string) error {
	caCert, err := ioutil.ReadFile(filepath.Join(CAPath, "ca.cert"))
	if err != nil {
		return err
	}

	caKey, err := ioutil.ReadFile(filepath.Join(CAPath, "ca.key"))
	if err != nil {
		return err
	}

	ca, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}

	if ca.Leaf, err = x509.ParseCertificate(ca.Certificate[0]); err != nil {
		return err
	}

	goproxy.GoproxyCa = ca
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&ca)}

	return nil
}

func (s *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}
