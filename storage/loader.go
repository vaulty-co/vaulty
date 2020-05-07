package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/vaulty/proxy/model"
)

type RoutesFile struct {
	Vault  *model.Vault `json:"vault"`
	Routes struct {
		Inbound  []*model.Route `json:"inbound"`
		Outbound []*model.Route `json:"outbound"`
	} `json:"routes"`
}

func LoadFromFile(file string, storage Storage) error {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	routesFile := &RoutesFile{}
	err = json.Unmarshal(fileContent, routesFile)
	if err != nil {
		return err
	}

	vault := routesFile.Vault
	storage.CreateVault(vault)

	for _, route := range routesFile.Routes.Inbound {
		route.Type = model.RouteInbound
		route.VaultID = vault.ID
		err = storage.CreateRoute(route)
		if err != nil {
			return err
		}
	}

	for _, route := range routesFile.Routes.Outbound {
		route.Type = model.RouteOutbound
		route.VaultID = vault.ID
		err = storage.CreateRoute(route)
		if err != nil {
			return err
		}
	}

	return nil
}
