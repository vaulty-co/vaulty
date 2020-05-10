package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/proxy/model"
)

type route struct {
	Method                  string
	Path                    string
	RequestTransformations  []map[string]interface{} `mapstructure:"request_transformations"`
	ResponseTransformations []map[string]interface{} `mapstructure:"response_transformations"`
}

type routesFile struct {
	Vault  map[string]interface{}
	Routes struct {
		Inbound  []*route
		Outbound []*route
	}
}

func LoadFromFile(file string, storage Storage) error {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var result routesFile
	err = json.Unmarshal(fileContent, &result)
	if err != nil {
		return err
	}

	var vault model.Vault
	mapstructure.Decode(result.Vault, &vault)
	storage.CreateVault(&vault)

	for _, rt := range result.Routes.Inbound {
		route := model.Route{
			Type:    model.RouteInbound,
			VaultID: vault.ID,
		}
		mapstructure.Decode(rt, &route)
		err = storage.CreateRoute(&route)
		if err != nil {
			return err
		}
	}

	for _, rt := range result.Routes.Outbound {
		route := model.Route{
			Type:    model.RouteOutbound,
			VaultID: vault.ID,
		}
		mapstructure.Decode(rt, &route)
		err = storage.CreateRoute(&route)
		if err != nil {
			return err
		}
	}

	return nil
}
