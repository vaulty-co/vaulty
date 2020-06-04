package transformer

import (
	"net/http"

	"github.com/vaulty/vaulty/action"
)

type RequestTransformer interface {
	TransformRequest(req *http.Request) (*http.Request, error)
}

type ResponseTransformer interface {
	TransformResponse(req *http.Response) (*http.Response, error)
}

type Transformer interface {
	RequestTransformer
	ResponseTransformer
}

// var TransformerRegistry map[string]Factory = map[string]Factory{}

type Factory func(map[string]interface{}, action.Action) (Transformer, error)

// func Builder(rawInput interface{}, act action.Action) (Transformer, error) {
// 	params := rawInput.(map[string]interface{})
// 	type_ := params["type"].(string)

// 	if factory, ok := TransformerRegistry[type_]; ok {
// 		transformation, err := factory(params, act)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return transformation, nil
// 	}

// 	return nil, errors.New(fmt.Sprintf("Unknown transformation type %s", type_))
// }
