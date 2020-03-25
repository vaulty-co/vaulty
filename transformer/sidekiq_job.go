package transformer

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type sidekiqJob struct {
	WorkerClass string        `json:"class"`
	Args        []interface{} `json:"args"`
	Retry       bool          `json:"retry"`
	JID         string        `json:"jid"`
	CreatedAt   int64         `json:"created_at"`
	EnqueuedAt  int64         `json:"enqueued_at"`
}

func newSidekiqJob(workerClass string, payload interface{}) *sidekiqJob {
	args := make([]interface{}, 1)
	args[0] = payload
	return &sidekiqJob{
		WorkerClass: workerClass,
		Args:        args,
		Retry:       false,
		JID:         genID(),
	}
}

func (job *sidekiqJob) JSON() ([]byte, error) {
	job.CreatedAt = time.Now().UnixNano()
	job.EnqueuedAt = time.Now().UnixNano()

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}
	return jobJSON, nil
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
