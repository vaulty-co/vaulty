package model

type Route struct {
	ID       string
	Upstream string
}

const (
	RouteInbound   = "INBOUND"
	RounteOutbound = "OUTBOUND"
)
