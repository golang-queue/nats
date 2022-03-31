package nats

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-queue/queue"

	"github.com/nats-io/nats.go"
)

var _ queue.Worker = (*Worker)(nil)

// Worker for NSQ
type Worker struct {
	client   *nats.Conn
	stop     chan struct{}
	stopFlag int32
	stopOnce sync.Once
	opts     options
}

// NewWorker for struc
func NewWorker(opts ...Option) *Worker {
	var err error
	w := &Worker{
		opts: newOptions(opts...),
		stop: make(chan struct{}),
	}

	w.client, err = nats.Connect(w.opts.addr)
	if err != nil {
		panic(err)
	}

	return w
}

func (w *Worker) handle(job queue.Job) error {
	// create channel with buffer size 1 to avoid goroutine leak
	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), job.Timeout)
	defer func() {
		cancel()
	}()

	// run the job
	go func() {
		// handle panic issue
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		// run custom process function
		done <- w.opts.runFunc(ctx, job)
	}()

	select {
	case p := <-panicChan:
		panic(p)
	case <-ctx.Done(): // timeout reached
		return ctx.Err()
	case <-w.stop: // shutdown service
		// cancel job
		cancel()

		leftTime := job.Timeout - time.Since(startTime)
		// wait job
		select {
		case <-time.After(leftTime):
			return context.DeadlineExceeded
		case err := <-done: // job finish
			return err
		case p := <-panicChan:
			panic(p)
		}
	case err := <-done: // job finish
		return err
	}
}

// Run start the worker
func (w *Worker) Run(task queue.QueuedMessage) error {
	wg := &sync.WaitGroup{}
	panicChan := make(chan interface{}, 1)
	_, err := w.client.QueueSubscribe(w.opts.subj, w.opts.queue, func(m *nats.Msg) {
		wg.Add(1)
		defer func() {
			wg.Done()
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		var data queue.Job
		_ = json.Unmarshal(m.Data, &data)

		if err := w.handle(data); err != nil {
			w.opts.logger.Error(err)
		}
	})
	if err != nil {
		return err
	}

	// wait close signal
	select {
	case <-w.stop:
	case err := <-panicChan:
		w.opts.logger.Error(err)
	}

	// wait job completed
	wg.Wait()

	return nil
}

// Shutdown worker
func (w *Worker) Shutdown() error {
	if !atomic.CompareAndSwapInt32(&w.stopFlag, 0, 1) {
		return queue.ErrQueueShutdown
	}

	w.stopOnce.Do(func() {
		w.client.Close()
		close(w.stop)
	})
	return nil
}

// Capacity for channel
func (w *Worker) Capacity() int {
	return 0
}

// Usage for count of channel usage
func (w *Worker) Usage() int {
	return 0
}

// Queue send notification to queue
func (w *Worker) Queue(job queue.QueuedMessage) error {
	if atomic.LoadInt32(&w.stopFlag) == 1 {
		return queue.ErrQueueShutdown
	}

	err := w.client.Publish(w.opts.subj, job.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Request a new task
func (w *Worker) Request() (queue.QueuedMessage, error) {
	return nil, nil
}
