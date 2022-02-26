package daemon

import (
	"context"
	"sync"
)

// Daemon orchestrates processing and sending incoming messages
type Daemon struct {
	Queue chan Message

	stop context.CancelFunc
	wg   *sync.WaitGroup
}

// New spawns a new sender daemon to process all the incoming messages
func New(workers int) *Daemon {
	// Allow gracefully stopping the daemon
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Create the incoming work queue
	queue := make(chan Message)

	// Spawn the workers
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(ctx, queue, &wg)
	}

	return &Daemon{
		Queue: queue,
		stop:  cancel,
		wg:    &wg,
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
