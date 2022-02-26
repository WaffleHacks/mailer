package daemon

import (
	"context"
	"fmt"
	"sync"

	gonanoid "github.com/matoous/go-nanoid"
	"go.uber.org/zap"
)

type BodyType int32

const (
	BodyTypePlain BodyType = iota + 1
	BodyTypeHTML
)

// Message represents an email to be sent to one or more recipients
type Message struct {
	To      []string
	From    string
	Subject string
	Body    string
	Type    BodyType
	ReplyTo *string
}

// worker processes and sends the incoming messages
func worker(ctx context.Context, queue <-chan Message, wg *sync.WaitGroup) {
	l := zap.L().Named("daemon:worker").With(zap.String("id", gonanoid.MustID(8)))
	l.Info("worker started")

	for {
		select {
		case message := <-queue:
			// TODO: actually send messages
			fmt.Printf("%+v\n", message)

		case <-ctx.Done():
			l.Info("worker exited")
			wg.Done()
		}
	}
}
