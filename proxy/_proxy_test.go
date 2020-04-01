package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var acceptAllCerts = &tls.Config{InsecureSkipVerify: true}

type EchoHandler struct{}

func (EchoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, readRequestBody(req))
}

var upstream = httptest.NewServer(EchoHandler{})

func setVaultUpstream(vaultID, upstream string) {
	upstreamKey := fmt.Sprintf("vault:%s:upstream", vaultID)
	redisClient.Set(upstreamKey, upstream, 0)
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

	req, _ := http.NewRequest(http.MethodPost, proxy.URL, nil)
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
