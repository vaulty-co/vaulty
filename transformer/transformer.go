package transformer

import (
	"net/http"

	"github.com/vaulty/proxy/model"
)

type Transformer interface {
	TransformRequestBody(route *model.Route, httpRequest *http.Request) error
	TransformResponseBody(routeID string, httpResponse *http.Response) error
}
