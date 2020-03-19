package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var acceptAllCerts = &tls.Config{InsecureSkipVerify: true}

type EchoHandler struct{}

func (EchoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, readRequestBody(req))
}

// func init() {
// 	http.DefaultServeMux.Handle("/", EchoHandler{})
// }

// var https = httptest.NewTLSServer(nil)
var upstream = httptest.NewServer(EchoHandler{})

func oneShotProxy(proxy *Proxy) (client *http.Client, s *httptest.Server) {
	s = httptest.NewServer(proxy.server)

	proxyUrl, _ := url.Parse(s.URL)
	tr := &http.Transport{TLSClientConfig: acceptAllCerts, Proxy: http.ProxyURL(proxyUrl)}
	client = &http.Client{Transport: tr}
	return
}

func setVaultUpstream(vaultID, upstream string) {
	_vault := &vault{
		ID: vaultID,
	}
	_vault.setUpstream(upstream)
}

func readRequestBody(req *http.Request) string {
	request, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return string(request)
}

func readResponseBody(res *http.Response) string {
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return string(response)
}

func TestIboundRequestForwardedToUpstream(t *testing.T) {
	setVaultUpstream("vlt123", upstream.URL)

	proxy := httptest.NewServer(NewProxy().server)
	defer proxy.Close()

	req, _ := http.NewRequest("POST", proxy.URL, nil)
	req.Host = "vlt123.proxy.test"

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	body := readResponseBody(res)
	fmt.Println(body)

}

// func TestNoVaultForIboundRequestToUpstream(t *testing.T) {
// 	proxy := httptest.NewServer(NewProxy().server)
// 	defer proxy.Close()

// 	req, _ := http.NewRequest("POST", proxy.URL, nil)
// 	req.Host = "vlt123.proxy.test"

// 	client := &http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	body := readResponseBody(res)
// 	fmt.Println(body)

// }

// func TestVaultNotFound(t *testing.T) {
// 	client, ts := oneShotProxy(NewProxy())
// 	defer ts.Close()

// 	// client := &http.Client{}
// 	req, _ := http.NewRequest("POST", ts.URL, nil)
// 	req.Host = "bla.proxy.test"
// 	res, err := client.Do(req)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if res.StatusCode != http.StatusNotFound {
// 		t.Errorf("Received %d http status, want %d", res.StatusCode, http.StatusNotFound)
// 	}

// 	// greeting, err := ioutil.ReadAll(res.Body)
// 	// res.Body.Close()
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// fmt.Printf("%s", greeting)
// }
