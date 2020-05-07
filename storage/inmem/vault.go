package inmem

import (
	"errors"

	"github.com/rs/xid"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
)

func (s *inmemStorage) CreateVault(vault *model.Vault) error {
	vault.ID = "vlt" + xid.New().String()

	s.vaults[vault.ID] = vault

	return nil
}

func (s *inmemStorage) ListVaults() ([]*model.Vault, error) {
	vaults := []*model.Vault{}

	for _, v := range s.vaults {
		vault := &model.Vault{}
		vault.ID = v.ID
		vault.Upstream = v.Upstream

		vaults = append(vaults, vault)
	}

	return vaults, nil
}

func (s *inmemStorage) FindVault(vaultID string) (*model.Vault, error) {
	if vaultID == "vltError" {
		return nil, errors.New("Test error")
	}

	vault, ok := s.vaults[vaultID]
	if !ok {
		// vault was not found
		return nil, storage.ErrNoRows
	}

	return vault, nil
}

func (s *inmemStorage) DeleteVault(vaultID string) error {
	err := s.DeleteRoutes(vaultID)
	if err != nil {
		return err
	}

	delete(s.vaults, vaultID)

	return nil
}

func (s *inmemStorage) UpdateVault(vault *model.Vault) error {
	s.vaults[vault.ID] = vault

	return nil
}
