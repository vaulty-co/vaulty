package transformer

import "net/http"

type Transformer interface {
	TransformRequestBody(routeID string, httpRequest *http.Request) error
	TransformResponseBody(routeID string, httpResponse *http.Response) error
}
