package model

import (
	"net/url"
)

type Vault struct {
	ID          string
	UpstreamURL *url.URL
}

func NewVault(id string, upstreamURL *url.URL) *Vault {
	return &Vault{
		ID:          id,
		UpstreamURL: upstreamURL,
	}
}

// func (v *Vault) GetUpstream() string {
// 	return v.redisClient.Get(fmt.Sprintf("vault:%s:upstream", v.ID)).Val()
// }

// func (v *Vault) setUpstream(upstream string) {
// 	v.redisClient.Set(fmt.Sprintf("vault:%s:upstream", v.ID), upstream, 0)
// }
