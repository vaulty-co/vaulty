package test_transformer

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vaulty/proxy/transformer"
)

type Transformer struct {
}

func NewTransformer() transformer.Transformer {
	return &Transformer{}
}

func (t *Transformer) TransformRequestBody(routeID string, req *http.Request) error {
	log.Printf("Transforming request body for route: %s", routeID)

	b, _ := ioutil.ReadAll(req.Body)
	body := string(b)

	body += " transformed"

	size := int64(len(body))

	req.Body = ioutil.NopCloser(bufio.NewReader(bytes.NewBufferString(body)))
	req.Header.Del("Content-Length")
	req.ContentLength = size

	log.Println("Done")

	return nil
}

func (t *Transformer) TransformResponseBody(routeID string, res *http.Response) error {
	b, _ := ioutil.ReadAll(res.Body)
	body := string(b)

	log.Printf("Transforming response body for route: %s", routeID)

	body += " transformed"

	res.Body = ioutil.NopCloser(bufio.NewReader(bytes.NewBufferString(body)))

	size := int64(len(body))
	res.Header.Del("Content-Length")
	res.ContentLength = size

	log.Println("Done")

	return nil
}
