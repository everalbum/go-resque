package resque

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type Job struct {
	Class string   `json:"class"`
	Args  []interface{} `json:"args"`
}

func NewJob(job_class string, args []interface{}) *Job {
	return &Job{job_class, makeJobArgs(args)}
}

func (j *Job) Encode() (jsonString string) {
	if jsonBytes, err := json.Marshal(&j); err == nil {
		jsonString = string(jsonBytes)
	}

	return
}

func Enqueue(client redis.Conn, queue, job_class string, args ...interface{}) (int64, error) {
	job := NewJob(job_class, args)

	return redis.Int64(client.Do("LPUSH", "resque:queue:"+queue, job.Encode()))
}

func makeJobArgs(args []interface{}) []interface{} {
	if len(args) == 0 {
		// NOTE: Dirty hack to make a [{}] JSON struct
		return append(make([]interface{}, 0), make(map[string]interface{}, 0))
	}

	return args
}
