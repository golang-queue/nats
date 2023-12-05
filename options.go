package nats

import (
	"context"
	"strings"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"

	"github.com/nats-io/nats.go"
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
func WithAddr(addrs ...string) Option {
	return func(w *options) {
		if len(addrs) > 0 {
			w.addr = strings.Join(addrs, ",")
		}
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
		addr:   nats.DefaultURL,
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
