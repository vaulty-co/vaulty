package transformer

import (
	"io/ioutil"
	"net/http"

	"github.com/vaulty/proxy/model"
)

// Serializable structure for http.Request and http.Response
// with raw body
type Request struct {
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	RouteID string      `json:"route_id"`
	VaultID string      `json:"vault_id"`
}

func newSerializableRequest(route *model.Route, req *http.Request) (*Request, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return &Request{
		Headers: req.Header,
		Body:    body,
		URL:     req.URL.String(),
		Method:  req.Method,
		RouteID: route.ID,
		VaultID: route.VaultID,
	}, nil
}
