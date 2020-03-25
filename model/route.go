package model

import "net/url"

type Route struct {
	ID          string
	UpstreamURL *url.URL
}

const (
	RouteInbound   = "inbound"
	RounteOutbound = "outbound"
)
