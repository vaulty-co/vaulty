package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/vaulty/model"
	"github.com/vaulty/vaulty/transform"
	"github.com/vaulty/vaulty/transform/action"
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

type LoaderOptions struct {
	ActionOptions *action.Options
	Storage       Storage
}

func LoadFromFile(file string, opts *LoaderOptions) error {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return Load(fileContent, opts)
}

func Load(rawJson []byte, opts *LoaderOptions) error {
	storage := opts.Storage
	actionOptions := opts.ActionOptions

	var input routesFile
	err := json.Unmarshal(rawJson, &input)
	if err != nil {
		return err
	}

	var vault model.Vault
	mapstructure.Decode(input.Vault, &vault)
	storage.CreateVault(&vault)

	inboundRoutes, err := buildRoutes(input.Routes.Inbound, vault.ID, model.RouteInbound, actionOptions)
	if err != nil {
		return err
	}

	outboundRoutes, err := buildRoutes(input.Routes.Outbound, vault.ID, model.RouteOutbound, actionOptions)
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

func buildRoutes(rawRoutes []*route, vaultID string, routesType model.RouteType, actionOptions *action.Options) ([]*model.Route, error) {
	var routes []*model.Route

	for _, rt := range rawRoutes {
		requestTransformations, err := buildTransformations(rt.RequestTransformations, actionOptions)
		if err != nil {
			return nil, err
		}

		responseTransformations, err := buildTransformations(rt.ResponseTransformations, actionOptions)
		if err != nil {
			return nil, err
		}

		route := model.Route{
			Type:                    routesType,
			VaultID:                 vaultID,
			RequestTransformations:  requestTransformations,
			ResponseTransformations: responseTransformations,
		}
		mapstructure.Decode(rt, &route)
		routes = append(routes, &route)
	}

	return routes, nil
}

func buildTransformations(rawTransformations []map[string]interface{}, actionOptions *action.Options) ([]transform.Transformer, error) {
	var transformations []transform.Transformer

	for _, tr := range rawTransformations {
		action, err := action.Factory(tr["action"], actionOptions)
		if err != nil {
			return nil, err
		}

		transformation, err := transform.Factory(tr, action)
		if err != nil {
			return nil, err
		}

		transformations = append(transformations, transformation)
	}

	return transformations, nil
}
