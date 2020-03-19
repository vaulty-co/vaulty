package proxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

func errResponse(r *http.Request, message string, status int) *http.Response {
	return goproxy.NewResponse(r,
		goproxy.ContentTypeText,
		status,
		message)
}
