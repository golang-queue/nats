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

// Option for queue system
type Option func(*Worker)

// Worker for NSQ
type Worker struct {
	addr        string
	subj        string
	queue       string
	client      *nats.Conn
	stop        chan struct{}
	stopOnce    sync.Once
	runFunc     func(context.Context, queue.QueuedMessage) error
	logger      queue.Logger
	stopFlag    int32
	busyWorkers uint64
}

func (w *Worker) incBusyWorker() {
	atomic.AddUint64(&w.busyWorkers, 1)
}

func (w *Worker) decBusyWorker() {
	atomic.AddUint64(&w.busyWorkers, ^uint64(0))
}

func (w *Worker) BusyWorkers() uint64 {
	return atomic.LoadUint64(&w.busyWorkers)
}

// WithAddr setup the addr of NATS
func WithAddr(addr string) Option {
	return func(w *Worker) {
		w.addr = "nats://" + addr
	}
}

// WithSubj setup the subject of NATS
func WithSubj(subj string) Option {
	return func(w *Worker) {
		w.subj = subj
	}
}

// WithQueue setup the queue of NATS
func WithQueue(queue string) Option {
	return func(w *Worker) {
		w.queue = queue
	}
}

// WithRunFunc setup the run func of queue
func WithRunFunc(fn func(context.Context, queue.QueuedMessage) error) Option {
	return func(w *Worker) {
		w.runFunc = fn
	}
}

// WithLogger set custom logger
func WithLogger(l queue.Logger) Option {
	return func(w *Worker) {
		w.logger = l
	}
}

// NewWorker for struc
func NewWorker(opts ...Option) *Worker {
	var err error
	w := &Worker{
		addr:  "127.0.0.1:4222",
		subj:  "foobar",
		queue: "foobar",
		stop:  make(chan struct{}),
		runFunc: func(context.Context, queue.QueuedMessage) error {
			return nil
		},
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(w)
	}

	w.client, err = nats.Connect(w.addr)
	if err != nil {
		panic(err)
	}

	return w
}

// BeforeRun run script before start worker
func (w *Worker) BeforeRun() error {
	return nil
}

// AfterRun run script after start worker
func (w *Worker) AfterRun() error {
	return nil
}

func (w *Worker) handle(job queue.Job) error {
	// create channel with buffer size 1 to avoid goroutine leak
	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), job.Timeout)
	w.incBusyWorker()
	defer func() {
		cancel()
		w.decBusyWorker()
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
		done <- w.runFunc(ctx, job)
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
func (w *Worker) Run() error {
	wg := &sync.WaitGroup{}
	panicChan := make(chan interface{}, 1)
	_, err := w.client.QueueSubscribe(w.subj, w.queue, func(m *nats.Msg) {
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
			w.logger.Error(err)
		}
	})
	if err != nil {
		return err
	}

	// wait close signal
	select {
	case <-w.stop:
	case err := <-panicChan:
		w.logger.Error(err)
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

	err := w.client.Publish(w.subj, job.Bytes())
	if err != nil {
		return err
	}

	return nil
}
