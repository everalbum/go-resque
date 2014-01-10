package resque

import (
  "encoding/json"
  "github.com/garyburd/redigo/redis"
)

type jobArg interface{}

type job struct {
  Class string   `json:"class"`
  Args  []jobArg `json:"args"`
}

func Enqueue(client *redis.Conn, queue, job_class string, args ...jobArg) int64, err {
  var j = &job{job_class, makeJobArgs(args)}

  job_json, _ := json.Marshal(j)

  return redis.Int64(client.Do("LPUSH", "resque:queue:"+queue, string(job_json[:])))
}

func makeJobArgs(args []jobArg) []jobArg {
  if len(args) == 0 {
    // NOTE: Dirty hack to make a [{}] JSON struct
    return append(make([]jobArg, 0), make(map[string]jobArg, 0))
  }

  return args
}
