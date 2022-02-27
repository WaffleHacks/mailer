package daemon

import (
	"context"
	"sync"
)

// Daemon orchestrates processing and sending incoming messages
type Daemon struct {
	Queue chan Message

	matchers []*Matcher
	stop     context.CancelFunc
	wg       *sync.WaitGroup
}

// New spawns a new sender daemon to process all the incoming messages
func New(matchers []*Matcher) *Daemon {
	// Allow gracefully stopping the daemon
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Create the incoming work queue
	queue := make(chan Message)

	// Spawn the workers
	for _, matcher := range matchers {
		wg.Add(matcher.workers)
		for i := 0; i < matcher.workers; i++ {
			go worker(ctx, matcher, &wg)
		}
	}

	d := &Daemon{
		matchers: matchers,
		Queue:    queue,
		stop:     cancel,
		wg:       &wg,
	}
	go d.dispatcher(ctx)

	return d
}

// dispatcher routes messages to the different providers
func (d *Daemon) dispatcher(ctx context.Context) {
	d.wg.Add(1)

	for {
		select {
		case msg := <-d.Queue:
			for _, matcher := range d.matchers {
				matcher.Enqueue(msg)
			}

		case <-ctx.Done():
			d.wg.Done()
			return
		}
	}
}

// Shutdown stops the processing of messages. Any unsent messages will be dropped.
func (d *Daemon) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		d.stop()
		d.wg.Wait()
		done <- struct{}{}
	}()

	// Race the context and done channels
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
