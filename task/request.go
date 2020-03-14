package task

import (
	"io/ioutil"
	"net/http"
)

// Serializable structure for http.Request and http.Response
// with raw body
type Request struct {
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
	URL     string      `json:"url"`
	Method  string      `json:"method"`
}

func NewRequest(req *http.Request) (*Request, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return &Request{
		Headers: req.Header,
		Body:    body,
		URL:     req.URL.String(),
		Method:  req.Method,
	}, nil
}
