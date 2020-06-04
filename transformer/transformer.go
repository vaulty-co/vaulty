package transformer

import (
	"net/http"

	"github.com/vaulty/vaulty/action"
)

type Transformer interface {
	TransformRequest(req *http.Request) (*http.Request, error)
	TransformResponse(req *http.Response) (*http.Response, error)
}

type Factory func(map[string]interface{}, action.Action) (Transformer, error)
