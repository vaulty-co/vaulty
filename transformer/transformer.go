package transformer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis"
)

type Transformer struct {
	redisClient *redis.Client
}

func NewTransformer(redisClient *redis.Client) *Transformer {
	// func NewRequestBodyTransformer(routeID string, httpRequest *http.Request) *Transformer {
	return &Transformer{
		redisClient: redisClient,
	}
}

func (t *Transformer) TransformRequestBody(routeID string, httpRequest *http.Request) error {
	request, err := newSerializableRequest(routeID, httpRequest)
	if err != nil {
		return err
	}

	result, err := t.transform("ProxyWorker::Worker", request)
	if err != nil {
		return err
	}

	body, size := newBody(result.Body)
	httpRequest.Header.Del("Content-Length")
	httpRequest.Body = body
	httpRequest.ContentLength = size

	return nil
}

func (t *Transformer) transform(workerClass string, payload interface{}) (*Response, error) {
	transformJob := newSidekiqJob(workerClass, payload)

	// TODO: we should add timeout here
	// to handle delays with sidekiq
	// Subscribe for job status
	pubsub := t.redisClient.Subscribe(transformJob.JID)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		return nil, err
	}

	ch := pubsub.ChannelSize(1)

	// Enqueue sidekiq job
	jobJSON, err := transformJob.JSON()

	if err != nil {
		return nil, err
	}

	_, err = t.redisClient.LPush("queue:default", jobJSON).Result()
	if err != nil {
		return nil, err
	}

	// Wait for task status
	status := <-ch
	if status.Payload != "done" {
		return nil, errors.New(fmt.Sprintf("Unexpected return from worker: %s", status.Payload))
	}

	rawResponse := t.redisClient.Get(transformJob.JID).Val()
	response := &Response{}
	err = json.Unmarshal([]byte(rawResponse), response)

	return response, err
}

func newBody(body []byte) (io.ReadCloser, int64) {
	bodyReader := bufio.NewReader(bytes.NewBuffer(body))
	size := int64(len(body))

	return ioutil.NopCloser(bodyReader), size
}
