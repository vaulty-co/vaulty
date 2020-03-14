package main

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
	"github.com/vaulty/proxy/task"
)

func transformBody(req *http.Request) error {
	request, err := task.NewRequest(req)
	if err != nil {
		return err
	}

	requestID := genRequestID()

	pubsub := redis.Client().Subscribe(requestID)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive()
	if err != nil {
		return err
	}

	ch := pubsub.ChannelSize(1)

	transformTask := task.NewTask("ProxyWorker::Worker", request, requestID)

	err = transformTask.Perform("default")
	if err != nil {
		return err
	}

	msg := <-ch
	if msg.Payload != "done" {
		return errors.New(fmt.Sprintf("Unexpected return from worker: %s", msg.Payload))
	}

	receivedPayload := redis.Client().Get(requestID).Val()
	response := &task.Response{}
	err = json.Unmarshal([]byte(receivedPayload), response)

	req.Header.Del("Content-Length")

	bodyBuff := bytes.NewBuffer(response.Body)
	req.ContentLength = int64(bodyBuff.Len())
	req.Body = ioutil.NopCloser(bufio.NewReader(bodyBuff))

	return nil
}

func genRequestID() string {
	// Return 12 random bytes as 24 character hex
	b := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
