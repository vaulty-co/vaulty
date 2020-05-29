package routing

import (
	"net/http"
	"net/url"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/vaulty/vaulty/transform"
)

type RouteParams struct {
	Name                    string
	Method                  string
	URL                     string
	Upstream                string
	RequestTransformations  []transform.Transformer
	ResponseTransformations []transform.Transformer
}

type Route struct {
	Name        string
	UpstreamURL *url.URL
	IsInbound   bool

	method                  string
	rawURL                  string
	url                     *url.URL
	requestTransformations  []transform.Transformer
	responseTransformations []transform.Transformer
}

func NewRoute(params *RouteParams) (*Route, error) {
	var err error
	route := &Route{
		Name:                    params.Name,
		method:                  params.Method,
		rawURL:                  params.URL,
		requestTransformations:  params.RequestTransformations,
		responseTransformations: params.ResponseTransformations,
	}

	route.url, err = url.Parse(params.URL)
	if err != nil {
		return nil, err
	}

	route.UpstreamURL, err = url.Parse(params.Upstream)
	if err != nil {
		return nil, err
	}

	route.IsInbound = !route.url.IsAbs()

	return route, nil
}

func (r *Route) Match(req *http.Request) bool {
	var matchingURL *url.URL

	// no need to do any checking for inbound request and outbound route
	if req.URL.Host == "inbound" && !r.IsInbound {
		return false
	}

	if req.URL.Host == "inbound" {
		matchingURL = &url.URL{}
		matchingURL.Path = req.URL.Path
	} else {
		matchingURL = &url.URL{}
		matchingURL.Scheme = req.URL.Scheme
		matchingURL.Host = req.URL.Host
		matchingURL.Path = req.URL.Path
	}

	if matchingURL.Path == "" {
		matchingURL.Path = "/"
	}

	// check if route URL matches request URL
	// here we use filepath.Match which seems to be pretty good
	// for our goal.
	urlMatch, err := filepath.Match(r.rawURL, matchingURL.String())
	if err != nil {
		log.Errorf("route URL has mailformed pattern: %s", r.rawURL)
		return false
	}

	return urlMatch && (r.method == "*" || req.Method == r.method)
}
