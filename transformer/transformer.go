package transformer

import (
	"net/http"

	"github.com/vaulty/proxy/model"
)

type Transformer interface {
	TransformRequestBody(route *model.Route, httpRequest *http.Request) error
	TransformResponseBody(route *model.Route, httpResponse *http.Response) error
}
