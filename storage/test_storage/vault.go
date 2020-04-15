package test_storage

import (
	"errors"

	"github.com/rs/xid"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
)

func (s *TestStorage) CreateVault(vault *model.Vault) error {
	vault.ID = "vlt" + xid.New().String()

	testVaults[vault.ID] = vault

	return nil
}

func (s *TestStorage) ListVaults() ([]*model.Vault, error) {
	vaults := []*model.Vault{}

	for _, v := range testVaults {
		vault := &model.Vault{}
		vault.ID = v.ID
		vault.Upstream = v.Upstream

		vaults = append(vaults, vault)
	}

	return vaults, nil
}

func (s *TestStorage) FindVault(vaultID string) (*model.Vault, error) {
	if vaultID == "vltError" {
		return nil, errors.New("Test error")
	}

	vault, ok := testVaults[vaultID]
	if !ok {
		// vault was not found
		return nil, storage.ErrNoRows
	}

	return vault, nil
}

func (s *TestStorage) DeleteVault(vaultID string) error {
	delete(testVaults, vaultID)

	return nil
}

func (s *TestStorage) UpdateVault(vault *model.Vault) error {
	testVaults[vault.ID] = vault

	return nil
}
