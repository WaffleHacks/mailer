package daemon

import (
	"context"
	"fmt"
	"sync"

	"github.com/k3a/html2text"
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
			plain, html := makeBodies(message.Body, message.Type)
			fmt.Println(plain, html)

			// TODO: actually send messages

		case <-ctx.Done():
			l.Info("worker exited")
			wg.Done()
			return
		}
	}
}

// makeBodies creates a plain text and, optionally, a HTML body based on the provided type
func makeBodies(content string, bodyType BodyType) (string, *string) {
	if bodyType == BodyTypePlain {
		return content, nil
	}

	return html2text.HTML2Text(content), &content
}
