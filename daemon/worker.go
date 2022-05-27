package daemon

import (
	"context"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/k3a/html2text"
	gonanoid "github.com/matoous/go-nanoid"
	"go.uber.org/zap"

	"github.com/WaffleHacks/mailer/logging"
	"github.com/WaffleHacks/mailer/providers"
)

type BodyType int32

const (
	BodyTypePlain BodyType = iota + 1
	BodyTypeHTML
)

// Message represents an email to be sent to one or more recipients
type Message struct {
	To       []string
	From     string
	Subject  string
	Body     string
	Type     BodyType
	ReplyTo  *string
	Template bool
	Context  []map[string]string
}

// worker processes and sends the incoming messages
func worker(ctx context.Context, matcher *Matcher, wg *sync.WaitGroup) {
	workerId := gonanoid.MustID(8)
	l := logging.L().Named("daemon:worker").With(zap.String("provider", matcher.id), zap.String("id", workerId))
	l.Info("worker started")

	batchedProvider, supportsBatching := matcher.provider.(providers.BatchedProvider)

	for {
		select {
		case message := <-matcher.queue:
			span := sentry.StartSpan(ctx, "message")
			span.SetTag("provider", matcher.id)
			span.SetTag("worker", workerId)

			plain, html := makeBodies(message.Body, message.Type)

			// Select batch or single sending
			var err error
			if len(message.To) == 1 {
				err = matcher.provider.Send(span.Context(), l, message.To[0], message.From, message.Subject, plain, html, message.ReplyTo)
			} else if supportsBatching {
				err = batchedProvider.SendBatch(span.Context(), l, message.To, message.From, message.Subject, plain, html, message.ReplyTo)
			} else {
				for _, to := range message.To {
					err = matcher.provider.Send(span.Context(), l, to, message.From, message.Subject, plain, html, message.ReplyTo)
					if err != nil {
						break
					}
				}
			}

			if err == nil {
				l.Info("sent message(s)", zap.Int("count", len(message.To)), zap.Strings("to", message.To))
			} else {
				l.Error("failed to send message(s)", zap.Error(err), zap.Strings("to", message.To))
			}
			span.Finish()

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
