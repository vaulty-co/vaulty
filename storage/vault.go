package storage

import (
	"github.com/rs/xid"
	"github.com/vaulty/proxy/model"
)

// We store vaults in the following way:
// vault:%s:upstream

func (s *redisStorage) CreateVault(vault *model.Vault) error {
	vault.ID = "vlt" + xid.New().String()

	err := s.redisClient.Set(vault.UpstreamKey(), vault.Upstream, 0).Err()
	if err != nil {
		return err
	}

	err = s.redisClient.LPush("vaults", vault.ID).Err()

	return err
}

func (s *redisStorage) ListVaults() ([]*model.Vault, error) {
	vaults := []*model.Vault{}

	res := s.redisClient.LRange("vaults", 0, -1)
	if res.Err() != nil {
		return nil, res.Err()
	}

	ids := res.Val()

	for _, id := range ids {
		vault := &model.Vault{}
		vault.ID = id
		vault.Upstream = s.redisClient.Get(vault.UpstreamKey()).Val()

		vaults = append(vaults, vault)
	}

	return vaults, nil
}

func (s *redisStorage) FindVault(vaultID string) (*model.Vault, error) {
	vault := &model.Vault{
		ID: vaultID,
	}

	vault.Upstream = s.redisClient.Get(vault.UpstreamKey()).Val()
	if vault.Upstream == "" {
		return nil, ErrNoRows
	}

	return vault, nil
}

func (s *redisStorage) UpdateVault(vault *model.Vault) error {
	err := s.redisClient.Set(vault.UpstreamKey(), vault.Upstream, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisStorage) DeleteVault(vaultID string) error {
	vault := &model.Vault{
		ID: vaultID,
	}

	err := s.redisClient.LRem("vaults", 1, vault.ID).Err()
	if err != nil {
		return err
	}

	return s.redisClient.Del(vault.UpstreamKey()).Err()
}
