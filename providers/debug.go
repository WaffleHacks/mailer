package providers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

var (
	debugTracer = otel.Tracer("github.com/WaffleHacks/mailer/providers/debug")

	debugSimulatedErrorAttr = attribute.Key("mailer.debug.simulated-error")
)

type Debug struct {
	failureRate int
	showBody    bool
	limiter     ratelimit.Limiter
}

func (d *Debug) shouldError() error {
	if rand.Intn(100) <= d.failureRate {
		return errors.New("simulated error")
	}

	return nil
}

func (d *Debug) formatMessage(to zap.Field, from, subject, body string, htmlBody, replyTo *string) []zap.Field {
	fields := []zap.Field{to, zap.String("from", from), zap.String("subject", subject), zap.Stringp("replyTo", replyTo), zap.Bool("isHtml", htmlBody != nil)}
	if d.showBody {
		fields = append(fields, zap.String("body", body), zap.Stringp("htmlBody", htmlBody))
	}

	return fields
}

func (d *Debug) Name() string {
	return "debug"
}

func (d *Debug) Send(ctx context.Context, l *zap.Logger, to, from, subject, body string, htmlBody, replyTo *string) error {
	_, span := debugTracer.Start(ctx, "send")
	defer span.End()

	d.limiter.Take()

	if err := d.shouldError(); err != nil {
		span.SetAttributes(debugSimulatedErrorAttr.Bool(true))
		l.Debug("simulating error")
		return err
	}

	l.Info("send email", d.formatMessage(zap.String("to", to), from, subject, body, htmlBody, replyTo)...)

	return nil
}

func (d *Debug) SendBatch(ctx context.Context, l *zap.Logger, to []string, from, subject, body string, htmlBody, replyTo *string) error {
	_, span := debugTracer.Start(ctx, "send-batch")
	defer span.End()

	if err := d.shouldError(); err != nil {
		span.SetAttributes(debugSimulatedErrorAttr.Bool(true))
		l.Debug("simulating error")
		return err
	}

	l.Info("send batch email", d.formatMessage(zap.Strings("to", to), from, subject, body, htmlBody, replyTo)...)

	return nil
}

// NewDebug creates a new Debug email provider
func NewDebug(id string) (Provider, error) {
	envId := strings.ToUpper(id)

	rawFailureRate := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_FAILURE_PERCENT", envId))
	failureRate := 0
	if len(rawFailureRate) != 0 {
		rate, err := strconv.Atoi(rawFailureRate)
		if err != nil {
			return nil, err
		} else if rate > 100 || rate < 0 {
			return nil, fmt.Errorf("%s: failure rate must be between 0 and 100, got %d", envId, rate)
		}

		failureRate = rate
	}

	rawRateLimit := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_RATE_LIMIT", envId))
	rl := ratelimit.NewUnlimited()
	if len(rawRateLimit) != 0 {
		perSecond, err := strconv.Atoi(rawRateLimit)
		if err != nil {
			return nil, err
		}

		if perSecond > 0 {
			rl = ratelimit.New(perSecond, ratelimit.Per(time.Second), ratelimit.WithSlack(10))
		}
	}

	rawShowBody := strings.ToLower(os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_SHOW_BODY", envId)))
	showBody := rawShowBody == "y" || rawShowBody == "yes" || rawShowBody == "t" || rawShowBody == "true"

	return &Debug{
		failureRate: failureRate,
		showBody:    showBody,
		limiter:     rl,
	}, nil
}
