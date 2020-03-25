package storage

import (
	"fmt"
	"net/url"

	"github.com/vaulty/proxy/model"
)

func (s *Storage) FindVault(vaultID string) (*model.Vault, error) {
	upstreamKey := fmt.Sprintf("vault:%s:upstream", vaultID)
	upstream := s.redisClient.Get(upstreamKey).Val()
	if upstream == "" {
		return nil, nil
	}

	upstreamURL, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}

	return &model.Vault{
		ID:          vaultID,
		UpstreamURL: upstreamURL,
	}, nil
}
