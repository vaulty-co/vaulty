package task

import (
	"encoding/json"
	"time"

	"github.com/vaulty/proxy/redis"
)

type Task struct {
	WorkerClass string      `json:"class"`
	Queue       string      `json:"queue"`
	Args        interface{} `json:"args"`
	Retry       bool        `json:"retry"`
	Jid         string      `json:"jid"`
	CreatedAt   int64       `json:"created_at"`
	EnqueuedAt  int64       `json:"enqueued_at"`
}

func NewTask(workerClass string, payload interface{}, jid string) *Task {
	return &Task{
		WorkerClass: workerClass,
		Args:        payload,
		Retry:       false,
		Jid:         jid,
	}
}

func (task *Task) Perform(queue string) error {
	task.CreatedAt = time.Now().UnixNano()
	task.EnqueuedAt = time.Now().UnixNano()
	task.Queue = queue

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	redis.Client().LPush("queue:"+queue, taskJSON)

	return nil
}
