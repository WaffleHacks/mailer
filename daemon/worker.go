package daemon

import (
	"context"
	"sync"

	"github.com/k3a/html2text"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/WaffleHacks/mailer/logging"
	"github.com/WaffleHacks/mailer/providers"
)

type BodyType int32

const (
	BodyTypePlain BodyType = iota + 1
	BodyTypeHTML
)

func (b BodyType) String() string {
	switch b {
	case BodyTypePlain:
		return "plain"
	case BodyTypeHTML:
		return "html"
	default:
		return "unknown"
	}
}

// Message represents an email to be sent to one or more recipients
type Message struct {
	To          []string
	From        string
	Subject     string
	Body        string
	Type        BodyType
	ReplyTo     *string
	SpanContext trace.SpanContext
}

var (
	tracer = otel.Tracer("github.com/WaffleHacks/mailer/daemon")

	toAttr       = attribute.Key("mailer.to")
	fromAttr     = attribute.Key("mailer.from")
	typeAttr     = attribute.Key("mailer.type")
	subjectAttr  = attribute.Key("mailer.subject")
	providerAttr = attribute.Key("mailer.worker.provider")
	workerAttr   = attribute.Key("mailer.worker.id")

	errorAttr            = attribute.Key("error")
	errorDescriptionAttr = attribute.Key("error.description")
)

// worker processes and sends the incoming messages
func worker(ctx context.Context, matcher *Matcher, wg *sync.WaitGroup) {
	workerId := gonanoid.Must(8)
	l := logging.L().Named("daemon:worker").With(zap.String("provider", matcher.id), zap.String("id", workerId))
	l.Info("worker started")

	batchedProvider, supportsBatching := matcher.provider.(providers.BatchedProvider)

	for {
		select {
		case message := <-matcher.queue:
			spanCtx, span := tracer.Start(ctx, "message",
				trace.WithAttributes(
					toAttr.StringSlice(message.To),
					fromAttr.String(message.From),
					typeAttr.String(message.Type.String()),
					subjectAttr.String(message.Subject),
					providerAttr.String(matcher.id),
					workerAttr.String(workerId),
				),
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithLinks(trace.Link{SpanContext: message.SpanContext}),
				trace.WithNewRoot(),
			)

			plain, html := makeBodies(spanCtx, message.Body, message.Type)

			// Select batch or single sending
			var err error
			if len(message.To) == 1 {
				err = matcher.provider.Send(spanCtx, l, message.To[0], message.From, message.Subject, plain, html, message.ReplyTo)
			} else if supportsBatching {
				err = batchedProvider.SendBatch(spanCtx, l, message.To, message.From, message.Subject, plain, html, message.ReplyTo)
			} else {
				for _, to := range message.To {
					err = matcher.provider.Send(spanCtx, l, to, message.From, message.Subject, plain, html, message.ReplyTo)
					if err != nil {
						break
					}
				}
			}

			if err == nil {
				l.Info("sent message(s)", zap.Int("count", len(message.To)), zap.Strings("to", message.To))
			} else {
				l.Error("failed to send message(s)", zap.Error(err), zap.Strings("to", message.To))
				span.SetAttributes(errorAttr.Bool(true), errorDescriptionAttr.String(err.Error()))
			}
			span.End()

		case <-ctx.Done():
			l.Info("worker exited")
			wg.Done()
			return
		}
	}
}

// makeBodies creates a plain text and, optionally, a HTML body based on the provided type
func makeBodies(ctx context.Context, content string, bodyType BodyType) (string, *string) {
	_, span := tracer.Start(ctx, "make-body")
	defer span.End()

	if bodyType == BodyTypePlain {
		return content, nil
	}

	return html2text.HTML2Text(content), &content
}
