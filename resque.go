package resque

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type Job struct {
	Class string   `json:"class"`
	Args  []interface{} `json:"args"`
}

func NewJob(jobClass string, args []interface{}) *Job {
	return &Job{jobClass, makeJobArgs(args)}
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

func Enqueue(client redis.Conn, queue, jobClass string, args ...interface{}) (int64, error) {
	job := NewJob(jobClass, args)

	return job.Enqueue(client, queue)
}

func makeJobArgs(args []interface{}) []interface{} {
	if len(args) == 0 {
		// NOTE: Dirty hack to make a [{}] JSON struct
		return append(make([]interface{}, 0), make(map[string]interface{}, 0))
	}

	return args
}
