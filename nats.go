package nats

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
	"github.com/golang-queue/queue/job"

	"github.com/nats-io/nats.go"
)

var _ core.Worker = (*Worker)(nil)

// Worker for NSQ
type Worker struct {
	client       *nats.Conn
	stop         chan struct{}
	exit         chan struct{}
	stopFlag     int32
	stopOnce     sync.Once
	startOnce    sync.Once
	opts         options
	subscription *nats.Subscription
	tasks        chan *nats.Msg
}

// NewWorker for struc
func NewWorker(opts ...Option) *Worker {
	var err error
	w := &Worker{
		opts:  newOptions(opts...),
		stop:  make(chan struct{}),
		exit:  make(chan struct{}),
		tasks: make(chan *nats.Msg),
	}

	w.client, err = nats.Connect(w.opts.addr)
	if err != nil {
		w.opts.logger.Fatal("can't connect to nats:", err)
	}

	if err := w.startConsumer(); err != nil {
		w.opts.logger.Fatal("can't start consumer:", err)
	}

	return w
}

func (w *Worker) startConsumer() (err error) {
	w.startOnce.Do(func() {
		w.subscription, err = w.client.QueueSubscribe(w.opts.subj, w.opts.queue, func(msg *nats.Msg) {
			select {
			case w.tasks <- msg:
			case <-w.stop:
				if msg != nil {
					// re-queue the task if worker has been shutdown.
					w.opts.logger.Info("re-queue the current task")
					if err := w.client.Publish(w.opts.subj, msg.Data); err != nil {
						w.opts.logger.Errorf("error to re-queue the current task: %s", err.Error())
					}
				}
				close(w.exit)
			}
		})
	})

	return err
}

// Run start the worker
func (w *Worker) Run(ctx context.Context, task core.QueuedMessage) error {
	return w.opts.runFunc(ctx, task)
}

// Shutdown worker
func (w *Worker) Shutdown() error {
	if !atomic.CompareAndSwapInt32(&w.stopFlag, 0, 1) {
		return queue.ErrQueueShutdown
	}

	w.stopOnce.Do(func() {
		// unsubscribe channel if start the consumer
		if w.subscription != nil {
			_ = w.subscription.Unsubscribe()
		}

		close(w.stop)
		select {
		case <-w.exit:
		case <-time.After(50 * time.Millisecond):
		}
		w.client.Close()
		close(w.tasks)
	})
	return nil
}

// Queue send notification to queue
func (w *Worker) Queue(job core.QueuedMessage) error {
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
func (w *Worker) Request() (core.QueuedMessage, error) {
	_ = w.startConsumer()
	clock := 0
loop:
	for {
		select {
		case task, ok := <-w.tasks:
			if !ok {
				return nil, queue.ErrQueueHasBeenClosed
			}
			var data job.Message
			_ = json.Unmarshal(task.Data, &data)
			return &data, nil
		case <-time.After(1 * time.Second):
			if clock == 5 {
				break loop
			}
			clock += 1
		}
	}

	return nil, queue.ErrNoTaskInQueue
}
