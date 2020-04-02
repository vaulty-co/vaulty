package model

import (
	"fmt"
	"net/url"
)

type Vault struct {
	ID          string
	Upstream    string
	upstreamURL *url.URL
}

func NewVault(id, upstream string) *Vault {
	return &Vault{
		ID:       id,
		Upstream: upstream,
	}
}

func (v *Vault) UpstreamKey() string {
	return fmt.Sprintf("vault:%s:upstream", v.ID)
}

func (v *Vault) UpstreamURL() *url.URL {
	if v.upstreamURL != nil {
		return v.upstreamURL

	}

	// ignore error (_) here as we should validate
	// upstream URL when we create it, not when
	// we use it
	v.upstreamURL, _ = url.Parse(v.Upstream)

	return v.upstreamURL
}
