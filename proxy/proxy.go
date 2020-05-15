package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/elazarl/goproxy"
	"github.com/vaulty/proxy/core"
	"github.com/vaulty/proxy/storage"
)

type Proxy struct {
	server  *goproxy.ProxyHttpServer
	storage storage.Storage
	config  *core.Configuration
}

func NewProxy(storage storage.Storage, config *core.Configuration) (*Proxy, error) {
	if config.ProxyPassword == "" {
		return nil, errors.New("Proxy password must be specified via config file or PROXY_PASS environment variable")
	}

	server := goproxy.NewProxyHttpServer()

	err := setupCA(config)
	if err != nil {
		return nil, err
	}

	proxy := &Proxy{
		server:  server,
		storage: storage,
		config:  config,
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

func setupCA(config *core.Configuration) error {
	caCert, err := ioutil.ReadFile(filepath.Join(config.CaPath, "ca.pem"))
	if err != nil {
		return err
	}

	caKey, err := ioutil.ReadFile(filepath.Join(config.CaPath, "ca.key"))
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
