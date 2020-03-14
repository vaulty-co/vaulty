package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/vaulty/proxy/task"
)

func transformBody(req *http.Request) error {
	request, err := task.NewRequest(req)
	if err != nil {
		return err
	}

	requestID := genRequestID()

	transformTask := task.NewTask("ProxyWorker::Worker", request, requestID)

	err = transformTask.Perform("default")
	if err != nil {
		return err
	}

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
