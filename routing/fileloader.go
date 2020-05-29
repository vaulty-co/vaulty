package routing

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/vaulty/encrypt"
	"github.com/vaulty/vaulty/secrets"
	"github.com/vaulty/vaulty/transform"
	"github.com/vaulty/vaulty/transform/action"
)

type routeDef struct {
	Name                    string
	Method                  string
	URL                     string
	Upstream                string
	RequestTransformations  []map[string]interface{} `json:"request_transformations"`
	ResponseTransformations []map[string]interface{} `json:"response_transformations"`
}

type fileDef struct {
	Options struct {
		DefaultUpstream string `json:"default_upstream"`
	}
	Routes []*routeDef
}

type fileLoader struct {
	enc            encrypt.Encrypter
	secretsStorage secrets.SecretsStorage
}

func NewFileLoader(enc encrypt.Encrypter, secretsStorage secrets.SecretsStorage) *fileLoader {
	return &fileLoader{
		enc:            enc,
		secretsStorage: secretsStorage,
	}
}

func (l *fileLoader) Load(filename string) ([]*Route, error) {
	rawJSON, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var input fileDef
	err = json.Unmarshal(rawJSON, &input)
	if err != nil {
		return nil, err
	}

	actionOptions := &action.Options{
		Encrypter:      l.enc,
		SecretsStorage: l.secretsStorage,
	}

	var routes []*Route

	for _, rd := range input.Routes {
		requestTransformations, err := buildTransformations(rd.RequestTransformations, actionOptions)
		if err != nil {
			return nil, err
		}

		responseTransformations, err := buildTransformations(rd.ResponseTransformations, actionOptions)
		if err != nil {
			return nil, err
		}

		routeParams := RouteParams{
			RequestTransformations:  requestTransformations,
			ResponseTransformations: responseTransformations,
		}

		mapstructure.Decode(rd, &routeParams)

		if routeParams.Upstream == "" {
			routeParams.Upstream = input.Options.DefaultUpstream
		}

		route, err := NewRoute(&routeParams)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
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
