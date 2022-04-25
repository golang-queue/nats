package nats

import (
	"context"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
)

// Option for queue system
type Option func(*options)

type options struct {
	runFunc func(context.Context, core.QueuedMessage) error
	logger  queue.Logger
	addr    string
	subj    string
	queue   string
}

// WithAddr setup the addr of NATS
func WithAddr(addr string) Option {
	return func(w *options) {
		w.addr = "nats://" + addr
	}
}

// WithSubj setup the subject of NATS
func WithSubj(subj string) Option {
	return func(w *options) {
		w.subj = subj
	}
}

// WithQueue setup the queue of NATS
func WithQueue(queue string) Option {
	return func(w *options) {
		w.queue = queue
	}
}

// WithRunFunc setup the run func of queue
func WithRunFunc(fn func(context.Context, core.QueuedMessage) error) Option {
	return func(w *options) {
		w.runFunc = fn
	}
}

// WithLogger set custom logger
func WithLogger(l queue.Logger) Option {
	return func(w *options) {
		w.logger = l
	}
}

func newOptions(opts ...Option) options {
	defaultOpts := options{
		addr:   "127.0.0.1:4222",
		subj:   "foobar",
		queue:  "foobar",
		logger: queue.NewLogger(),
		runFunc: func(context.Context, core.QueuedMessage) error {
			return nil
		},
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(&defaultOpts)
	}

	return defaultOpts
}
