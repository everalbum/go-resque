# go-resque

Simple [Resque](https://github.com/defunkt/resque) queue client for [Go](http://golang.org).

## Installation

```
go get github.com/everalbum/go-resque
```

## Usage

For example, a simple Resque job being processed by a Ruby worker somewhere:

```ruby
module Demo
  class Job
    def self.perform(param=nil)
      puts "Processed a job! Param: #{param.inspect}"
    end
  end
end
```

Enqueue this job from Go:

```go
package main

import (
  "github.com/everalbum/go-resque"
  "github.com/garyburd/redigo/redis"
)

func main() {
  conn, err := redis.Dial("tcp", "127.0.0.1:6379") // Create new Redis client to use for enqueuing

  if err != nil {
    // Handle error
  }

  defer conn.Close()

  // Enqueue the job into the "email" queue with appropriate client
  resque.Enqueue(conn, "email", "Demo::Job")

  // Enqueue into the "default" queue with passing one parameter to the Demo::Job.perform
  resque.Enqueue(conn, "default", "Demo::Job", 1)

  // Enqueue into the "default" queue with passing multiple
  // parameters to Demo::Job.perform so it will fail
  resque.Enqueue(conn, "default", "Demo::Job", 1, 2, "woot")
}
```

This also works with [resque-scheduler](https://github.com/resque/resque-scheduler). You can enqueue jobs using `EnqueueIn` or `EnqueueAt` to enqueue a job in the future.

```
  // Enqueues this job 60 seconds from now.
  delay := time.Duration(60) * time.Second
  resque.EnqueueIn(conn, delay, "default", "Demo::Job", 1, 2, "woot")
  
  // Enqueues this job at a specific time
  t := time.Now().Add(delay)
  resque.EnqueueAt(conn, t, "default", "Demo::Job", 1, 2, "woot")
```
