package transformer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/vaulty/proxy/model"
)

type SidekiqTransformer struct {
	redisClient *redis.Client
}

const defaultTransformationTimeout = 2 * time.Second

func NewSidekiqTransformer(redisClient *redis.Client) Transformer {
	// func NewRequestBodyTransformer(routeID string, httpRequest *http.Request) *Transformer {
	return &SidekiqTransformer{
		redisClient: redisClient,
	}
}

func (t *SidekiqTransformer) TransformRequestBody(route *model.Route, httpRequest *http.Request) error {
	request, err := newSerializableRequest(route, httpRequest)
	if err != nil {
		return err
	}

	result, err := t.transform("ProxyWorker::RequestBodyTransformer", request)
	if err != nil {
		return err
	}

	body, size := newBody(result.Body)
	httpRequest.Header.Del("Content-Length")
	httpRequest.Body = body
	httpRequest.ContentLength = size

	return nil
}

func (t *SidekiqTransformer) TransformResponseBody(routeID string, httpResponse *http.Response) error {
	return nil
}

func (t *SidekiqTransformer) transform(workerClass string, payload interface{}) (*Result, error) {
	transformJob := newSidekiqJob(workerClass, payload)

	pubsub := t.redisClient.Subscribe("task_" + transformJob.JID)
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

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(defaultTransformationTimeout)
		timeout <- true
	}()

	// Wait for task status
	select {
	case status := <-ch:
		fmt.Println("status", status)

		if status.Payload != "done" {
			return nil, fmt.Errorf("Unexpected return from worker: %s", status.Payload)
		}
	case <-timeout:
		return nil, fmt.Errorf("Transformation task timeout: %s", transformJob.JID)
	}

	rawResult := t.redisClient.Get(transformJob.JID).Val()
	result := &Result{}
	err = json.Unmarshal([]byte(rawResult), result)

	return result, err
}

func newBody(body []byte) (io.ReadCloser, int64) {
	bodyReader := bufio.NewReader(bytes.NewBuffer(body))
	size := int64(len(body))

	return ioutil.NopCloser(bodyReader), size
}
