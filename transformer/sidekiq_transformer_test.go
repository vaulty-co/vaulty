package transformer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vaulty/proxy/model"
	"github.com/vaulty/proxy/storage"
)

func TestRequestTransformation(t *testing.T) {
	assert := assert.New(t)

	rs := storage.NewRedisStorage(redisClient)

	//create test route with transformations
	requestTransformationJSON := `
	[
	  {
	    "type":"json",
	    "expression":"$.number",
	    "action":{
	      "type":"conceal"
	    }
	  },
	  {
	    "type":"json",
	    "expression":"$.email",
	    "action":{
	      "type":"conceal",
	      "format":"email"
	    }
	  }
	]
	`

	route := &model.Route{
		ID:                         "rt1",
		Type:                       model.RouteInbound,
		Method:                     http.MethodPost,
		Path:                       "/tokenize",
		VaultID:                    "vlt1",
		Upstream:                   "http://example.com",
		RequestTransformationsJSON: requestTransformationJSON,
	}

	err := rs.CreateRoute(route)
	assert.NoError(err)

	tr := NewSidekiqTransformer(redisClient)

	jsonBody := `
	{
	  "email":"john@example.com",
	  "number":"4242424242424242"
	}
	`

	req, _ := http.NewRequest(http.MethodPost, "/tokenize", bytes.NewBufferString(jsonBody))
	req.Header["Content-Type"] = []string{"application/json"}

	err = tr.TransformRequestBody(route, req)
	assert.NoError(err)

	jsonBlob, err := ioutil.ReadAll(req.Body)

	jsonResult := make(map[string]string)

	err = json.Unmarshal(jsonBlob, &jsonResult)

	assert.NoError(err)
	assert.Regexp(`\w+@example.com`, jsonResult["email"])
	assert.Regexp(`[\w]{8}(-[\w]{4}){3}-[\w]{12}`, jsonResult["number"])
}
