package transformer

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestTransformRequestBody(t *testing.T) {
	transformer := NewTransformer(redisClient)

	body := bytes.NewBufferString("Hello")
	req := httptest.NewRequest("GET", "http://example.com/foo", body)

	go func() {
		err := transformer.TransformRequestBody("rt123", req)
		if err != nil {
			t.Error(err)
		}
	}()

	for {
		channels, err := redisClient.PubSubChannels("task_").Result()
		if err != nil {
			t.Error(err)
		}

		if len(channels) > 0
	}

	// body, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println(resp.StatusCode)
	// fmt.Println(resp.Header.Get("Content-Type"))
	// fmt.Println(string(body))
}
