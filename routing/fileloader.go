package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"github.com/vaulty/vaulty/action"
	"github.com/vaulty/vaulty/encryption"
	"github.com/vaulty/vaulty/secrets"
	"github.com/vaulty/vaulty/transformer"
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
	enc                encryption.Encrypter
	secretsStorage     secrets.Storage
	salt               string
	transformerFactory map[string]transformer.Factory
}

type FileLoaderOptions struct {
	Enc                encryption.Encrypter
	SecretsStorage     secrets.Storage
	Salt               string
	TransformerFactory map[string]transformer.Factory
}

func NewFileLoader(opts *FileLoaderOptions) *fileLoader {
	return &fileLoader{
		enc:                opts.Enc,
		secretsStorage:     opts.SecretsStorage,
		salt:               opts.Salt,
		transformerFactory: opts.TransformerFactory,
	}
}

func (l *fileLoader) Load(filename string) ([]*Route, error) {
	rawJSON, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to load routes file (%s); %s", filename, err)
	}

	var input fileDef
	err = json.Unmarshal(rawJSON, &input)
	if err != nil {
		return nil, err
	}

	actionOptions := &action.Options{
		Encrypter:      l.enc,
		SecretsStorage: l.secretsStorage,
		Salt:           l.salt,
	}

	var routes []*Route

	for _, rd := range input.Routes {
		requestTransformations, err := l.buildTransformations(rd.RequestTransformations, actionOptions)
		if err != nil {
			return nil, err
		}

		responseTransformations, err := l.buildTransformations(rd.ResponseTransformations, actionOptions)
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

func (l *fileLoader) buildTransformations(rawTransformations []map[string]interface{}, actionOptions *action.Options) ([]transformer.Transformer, error) {
	var transformations []transformer.Transformer

	for _, tr := range rawTransformations {
		action, err := action.Factory(tr["action"], actionOptions)
		if err != nil {
			return nil, err
		}

		type_ := tr["type"].(string)
		factory, ok := l.transformerFactory[type_]
		if !ok {
			return nil, fmt.Errorf(`Factory for transformation type "%s" was not found`, type_)
		}

		transformation, err := factory(tr, action)
		if err != nil {
			return nil, err
		}

		transformations = append(transformations, transformation)
	}

	return transformations, nil
}
