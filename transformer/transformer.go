package transformer

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/vaulty/proxy/redis"
)

type Transformer struct {
	httpRequest  *http.Request
	httpResponse *http.Response
	jobID        string
}

func NewRequestBodyTransformer(httpRequest *http.Request) *Transformer {
	return &Transformer{
		httpRequest: httpRequest,
		jobID:       genID(),
	}
}

func (t *Transformer) TransformRequestBody() error {
	request, err := newSerializableRequest(t.httpRequest)
	if err != nil {
		return err
	}

	result, err := t.transform("ProxyWorker::Worker", request)
	if err != nil {
		return err
	}

	body, size := newBody(result.Body)
	t.httpRequest.Header.Del("Content-Length")
	t.httpRequest.Body = body
	t.httpRequest.ContentLength = size

	return nil
}

func (t *Transformer) transform(workerClass string, payload interface{}) (*Response, error) {
	// Subscribe for job status
	pubsub := redis.Client().Subscribe(t.jobID)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		return nil, err
	}

	ch := pubsub.ChannelSize(1)

	// Enqueue sidekiq job
	transformJob := NewJob(workerClass, payload, t.jobID)
	err = transformJob.Perform("default")
	if err != nil {
		return nil, err
	}

	// Wait for task status
	status := <-ch
	if status.Payload != "done" {
		return nil, errors.New(fmt.Sprintf("Unexpected return from worker: %s", status.Payload))
	}

	rawResponse := redis.Client().Get(t.jobID).Val()
	response := &Response{}
	err = json.Unmarshal([]byte(rawResponse), response)

	return response, err
}

func newBody(body []byte) (io.ReadCloser, int64) {
	bodyReader := bufio.NewReader(bytes.NewBuffer(body))
	size := int64(len(body))

	return ioutil.NopCloser(bodyReader), size
}

func genID() string {
	// Return 12 random bytes as 24 character hex
	b := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
