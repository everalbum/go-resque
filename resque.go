package resque

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Job struct {
	Class string        `json:"class"`
	Args  []interface{} `json:"args"`
	Queue string        `json:"queue,omitempty"`
}

type Encoder interface {
	Encode() string
}

func NewJob(jobClass string, args []interface{}, queue string) *Job {
	return &Job{jobClass, makeJobArgs(args), queue}
}

func (j *Job) Encode() (jsonString string) {
	if jsonBytes, err := json.Marshal(&j); err == nil {
		jsonString = string(jsonBytes)
	}

	return
}

func (j *Job) Enqueue(client redis.Conn, queue string) (int64, error) {
	return redis.Int64(client.Do("LPUSH", "resque:queue:"+queue, j.Encode()))
}

func (j *Job) EnqueueAt(client redis.Conn, t time.Time, queue string) error {
	jsonString := j.Encode()

	queueKey := fmt.Sprintf("resque:delayed:%d", t.Unix())
	timestampsValue := fmt.Sprintf("delayed:%d", t.Unix())

	client.Send("MULTI")
	client.Send("RPUSH", queueKey, jsonString)
	client.Send("SADD", "resque:timestamps:"+jsonString, timestampsValue)
	client.Send("ZADD", "resque:delayed_queue_schedule", t.Unix(), t.Unix())
	_, err := client.Do("EXEC")

	return err
}

func Enqueue(client redis.Conn, queue, jobClass string, args ...interface{}) (int64, error) {
	job := NewJob(jobClass, args, "")

	return job.Enqueue(client, queue)
}

func EnqueueIn(client redis.Conn, seconds int, queue, jobClass string, args ...interface{}) error {
	job := NewJob(jobClass, args, queue)
	delay := time.Duration(seconds) * time.Second
	enqueueTime := time.Now().Add(delay)

	return job.EnqueueAt(client, enqueueTime, queue)
}

func makeJobArgs(args []interface{}) []interface{} {
	if len(args) == 0 {
		// NOTE: Dirty hack to make a [{}] JSON struct
		return append(make([]interface{}, 0), make(map[string]interface{}, 0))
	}

	return args
}
