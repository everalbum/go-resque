package resque

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
)

var realjson = `{"class":"test","args":[1,true,"hello",3.14]}`
var args = []interface{}{1, true, "hello", 3.14}

func equal(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func TestNewJob(t *testing.T) {
	job := newJob("", "test", args)

	if job.Class != "test" {
		t.Errorf("got class %q, want %q", job.Class, "test")
	}
	if !equal(job.Args, args) {
		t.Errorf("got args %q, want %q", job.Args, args)
	}
}

func TestJob_Encode(t *testing.T) {
	job := newJob("", "test", args)

	if json := job.encode(); json != realjson {
		t.Errorf("got encoded %q, want %q", job.encode(), realjson)
	}
}

type testConn struct {
	redis.Conn
}
