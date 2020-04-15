package model

import (
	"fmt"
	"net/url"
)

type Route struct {
	ID                          string
	Type                        RouteType
	Method                      string
	Path                        string
	VaultID                     string
	Upstream                    string
	RequestTransformationsJSON  string
	ResponseTransformationsJSON string
}

type RouteType string

const (
	RouteInbound  RouteType = "inbound"
	RouteOutbound RouteType = "outbound"
)

func (r *Route) IDKey() string {
	return fmt.Sprintf("vault:%s:route:%s", r.VaultID, r.ID)
}

// vlt2uYBrnYkUnEF:INBOUND:POST:/records => routeID
func (r *Route) Key() string {
	return fmt.Sprintf("%s:%s:%s:%s", r.VaultID, r.Type, r.Method, r.Path)
}

func (r *Route) RequestKey() string {
	return fmt.Sprintf("%s:%s:%s:%s", r.VaultID, r.Type, r.Method, r.Path)
}

func (r *Route) UpstreamKey() string {
	return fmt.Sprintf("route:%s:upstream", r.ID)
}

func (r *Route) UpstreamURL() *url.URL {
	u, _ := url.Parse(r.Upstream)
	// ignore error here as we should validate
	// upstream URL when we create it, not when
	// we use it

	return u
}

func (rt RouteType) MarshalBinary() ([]byte, error) {
	return []byte(rt), nil
}

func (rt RouteType) UnmarshalBinary(b []byte) error {
	rt = RouteType(b)

	return nil
}
