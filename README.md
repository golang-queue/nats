# nats

[![CodeQL](https://github.com/golang-queue/nats/actions/workflows/codeql.yaml/badge.svg)](https://github.com/golang-queue/nats/actions/workflows/codeql.yaml)
[![Run Testing](https://github.com/golang-queue/nats/actions/workflows/go.yml/badge.svg)](https://github.com/golang-queue/nats/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/golang-queue/nats/branch/main/graph/badge.svg?token=FFZN8E2ZZB)](https://codecov.io/gh/golang-queue/nats)
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-queue/nats)](https://goreportcard.com/report/github.com/golang-queue/nats)

NATS as backend with [Queue package](https://github.com/golang-queue/queue) (Connective Technology for Adaptive Edge & Distributed Systems)

## Testing

setup the nats server

```sh
docker run -d --name nats-main -p 4222:4222 -p 8222:8222 nats:latest
```

run the test

```sh
go test -v ./...
```

## Example

```go
package main

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "time"

  "github.com/golang-queue/nats"
  "github.com/golang-queue/queue"
)

type job struct {
  Message string
}

func (j *job) Bytes() []byte {
  b, err := json.Marshal(j)
  if err != nil {
    panic(err)
  }
  return b
}

func main() {
  taskN := 100
  rets := make(chan string, taskN)

  // define the worker
  w := nats.NewWorker(
    nats.WithAddr("127.0.0.1:4222"),
    nats.WithSubj("example"),
    nats.WithQueue("foobar"),
    nats.WithRunFunc(func(ctx context.Context, m queue.QueuedMessage) error {
      var v *job
      if err := json.Unmarshal(m.Bytes(), &v); err != nil {
        return err
      }
      rets <- v.Message
      return nil
    }),
  )

  // define the queue
  q, err := queue.NewQueue(
    queue.WithWorkerCount(10),
    queue.WithWorker(w),
  )
  if err != nil {
    log.Fatal(err)
  }

  // start the five worker
  q.Start()

  // assign tasks in queue
  for i := 0; i < taskN; i++ {
    go func(i int) {
      if err := q.Queue(&job{
        Message: fmt.Sprintf("handle the job: %d", i+1),
      }); err != nil {
        log.Fatal(err)
      }
    }(i)
  }

  // wait until all tasks done
  for i := 0; i < taskN; i++ {
    fmt.Println("message:", <-rets)
    time.Sleep(50 * time.Millisecond)
  }

  // shutdown the service and notify all the worker
  q.Release()
}
```
