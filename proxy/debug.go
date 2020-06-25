package proxy

import (
	"mime"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

func debugRequest(req *http.Request) {
	if !log.IsLevelEnabled(log.DebugLevel) {
		return
	}

	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
		contentType := req.Header.Get("Content-Type")
		switch mediaType, _, _ := mime.ParseMediaType(contentType); mediaType {
		case "application/json", "multipart/form-data", "application/x-www-form-urlencoded", "plain/text":
			buf, _ := httputil.DumpRequest(req, true)
			log.Debugf("Request Dump:\n%s\n", buf)
		default:
			buf, _ := httputil.DumpRequest(req, false)
			log.Debugf("Request Dump without body:\n%s\n", buf)
		}

	}
}

func debugResponse(res *http.Response) {
	if !log.IsLevelEnabled(log.DebugLevel) {
		return
	}

	switch res.Request.Method {
	case "GET", "POST", "PUT", "PATCH":
		contentType := res.Header.Get("Content-Type")
		switch mediaType, _, _ := mime.ParseMediaType(contentType); mediaType {
		case "application/json", "multipart/form-data", "application/x-www-form-urlencoded", "plain/text":
			buf, _ := httputil.DumpResponse(res, true)
			log.Debugf("Response Dump:\n%s\n", buf)
		default:
			buf, _ := httputil.DumpResponse(res, false)
			log.Debugf("Response Dump without body:\n%s\n", buf)
		}
	}
}
