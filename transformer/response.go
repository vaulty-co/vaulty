package transformer

import (
	"io/ioutil"
	"net/http"

	"github.com/vaulty/proxy/model"
)

// Serializable structure for Sidekiq Worker
// with raw body
type Response struct {
	Body    []byte `json:"body"`
	RouteID string `json:"route_id"`
	VaultID string `json:"vault_id"`
}

func newSerializableResponse(route *model.Route, res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Body:    body,
		RouteID: route.ID,
		VaultID: route.VaultID,
	}, nil
}
