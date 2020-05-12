package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/transform"
	"github.com/vaulty/proxy/transform/action"
)

type route struct {
	Method                  string
	Path                    string
	RequestTransformations  []map[string]interface{} `json:"request_transformations"`
	ResponseTransformations []map[string]interface{} `json:"response_transformations"`
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

	var input routesFile
	err = json.Unmarshal(fileContent, &input)
	if err != nil {
		return err
	}

	var vault model.Vault
	mapstructure.Decode(input.Vault, &vault)
	storage.CreateVault(&vault)

	inboundRoutes, err := buildRoutes(input.Routes.Inbound, vault.ID, model.RouteInbound)
	if err != nil {
		return err
	}

	outboundRoutes, err := buildRoutes(input.Routes.Inbound, vault.ID, model.RouteOutbound)
	if err != nil {
		return err
	}

	routes := append(inboundRoutes, outboundRoutes...)

	for _, rt := range routes {
		err = storage.CreateRoute(rt)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildRoutes(rawRoutes []*route, vaultID string, routesType model.RouteType) ([]*model.Route, error) {
	var routes []*model.Route

	for _, rt := range rawRoutes {
		requestTransformations, err := buildTransformations(rt.RequestTransformations)
		if err != nil {
			return nil, err
		}

		responseTransformations, err := buildTransformations(rt.ResponseTransformations)
		if err != nil {
			return nil, err
		}

		route := model.Route{
			Type:                    model.RouteInbound,
			VaultID:                 vaultID,
			RequestTransformations:  requestTransformations,
			ResponseTransformations: responseTransformations,
		}
		mapstructure.Decode(rt, &route)
		routes = append(routes, &route)
	}

	return routes, nil
}

func buildTransformations(rawTransformations []map[string]interface{}) ([]transform.Transformer, error) {
	var transformations []transform.Transformer

	for _, tr := range rawTransformations {
		var transformation transform.Transformer

		action, err := buildAction(tr["action"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		switch tr["type"] {
		case "json":
			jsonTransformation := &transform.Json{
				Action: action,
			}
			err := mapstructure.Decode(tr, jsonTransformation)
			if err != nil {
				return nil, err
			}
			transformation = jsonTransformation
		}

		transformations = append(transformations, transformation)
	}

	return transformations, nil
}

func buildAction(rawAction map[string]interface{}) (transform.Transformer, error) {
	switch rawAction["type"] {
	case "encrypt":
		result := &action.Encrypt{}
		err := mapstructure.Decode(rawAction, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "tokenize":
		result := &action.Tokenize{}
		err := mapstructure.Decode(rawAction, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown action type %s", rawAction["type"]))
	}
}
