package resque

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type job struct {
	Queue string        `json:"queue,omitempty"`
	Class string        `json:"class"`
	Args  []interface{} `json:"args"`
}

func newJob(queue, jobClass string, args []interface{}) job {
	return job{queue, jobClass, makeJobArgs(args)}
}

func (j job) encode() (jsonString string) {
	if jsonBytes, err := json.Marshal(&j); err == nil {
		jsonString = string(jsonBytes)
	}

	return
}

func (j job) enqueue(client redis.Conn, queue string) (int64, error) {
	return redis.Int64(client.Do("RPUSH", "resque:queue:"+queue, j.encode()))
}

func (j job) enqueueAt(client redis.Conn, t time.Time, queue string) error {
	jsonString := j.encode()

	queueKey := fmt.Sprintf("resque:delayed:%d", t.Unix())
	timestampsValue := fmt.Sprintf("delayed:%d", t.Unix())

	client.Send("MULTI")
	client.Send("RPUSH", queueKey, jsonString)
	client.Send("SADD", "resque:timestamps:"+jsonString, timestampsValue)
	client.Send("ZADD", "resque:delayed_queue_schedule", t.Unix(), t.Unix())
	_, err := client.Do("EXEC")

	return err
}

func (j job) failure(client redis.Conn, queue string) (int64, error) {
	return redis.Int64(client.Do("RPUSH", "resque:failed:"+queue, j.encode()))
}

func Enqueue(client redis.Conn, queue, jobClass string, args ...interface{}) (int64, error) {
	job := newJob(queue, jobClass, args)

	return job.enqueue(client, queue)
}

func EnqueueIn(client redis.Conn, delay time.Duration, queue, jobClass string, args ...interface{}) error {
	job := newJob(queue, jobClass, args)
	enqueueTime := time.Now().Add(delay)

	return job.enqueueAt(client, enqueueTime, queue)
}

func EnqueueAt(client redis.Conn, t time.Time, queue, jobClass string, args ...interface{}) error {
	job := newJob(queue, jobClass, args)

	return job.enqueueAt(client, t, queue)
}

func Failure(client redis.Conn, queue, jobClass string, args ...interface{}) (int64, error) {
	job := newJob(queue, jobClass, args)

	return job.failure(client, queue)
}

func Size(client redis.Conn, queue string) (int64, error) {
	return redis.Int64(client.Do("LLEN", "resque:queue:"+queue))
}

func makeJobArgs(args []interface{}) []interface{} {
	if len(args) == 0 {
		// NOTE: Dirty hack to make a [{}] JSON struct
		return append(make([]interface{}, 0), make(map[string]interface{}, 0))
	}

	return args
}
