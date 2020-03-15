package transformer

import (
	"encoding/json"
	"time"

	"github.com/vaulty/proxy/redis"
)

type Job struct {
	WorkerClass string      `json:"class"`
	Queue       string      `json:"queue"`
	Args        interface{} `json:"args"`
	Retry       bool        `json:"retry"`
	Jid         string      `json:"jid"`
	CreatedAt   int64       `json:"created_at"`
	EnqueuedAt  int64       `json:"enqueued_at"`
}

func NewJob(workerClass string, payload interface{}, jid string) *Job {
	return &Job{
		WorkerClass: workerClass,
		Args:        payload,
		Retry:       false,
		Jid:         jid,
	}
}

func (job *Job) Perform(queue string) error {
	job.CreatedAt = time.Now().UnixNano()
	job.EnqueuedAt = time.Now().UnixNano()
	job.Queue = queue

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = redis.Client().LPush("queue:"+queue, jobJSON).Result()
	return err
}
