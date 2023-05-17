package providers

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/otel"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

var sesTracer = otel.Tracer("github.com/WaffleHacks/mailer/providers/ses")

func newContent(data string) *types.Content {
	return &types.Content{
		Data:    &data,
		Charset: aws.String("UTF-8"),
	}
}

type SES struct {
	client *sesv2.Client
	rl     ratelimit.Limiter
}

func (s *SES) Name() string {
	return "ses"
}

func (s *SES) Send(ctx context.Context, _ *zap.Logger, to, from, subject, body string, htmlBody, replyTo *string) error {
	_, span := sesTracer.Start(ctx, "send")
	defer span.End()

	s.rl.Take()

	input := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: newContent(body),
				},
				Subject: newContent(subject),
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		FromEmailAddress: &from,
	}
	if htmlBody != nil {
		input.Content.Simple.Body.Html = newContent(*htmlBody)
	}
	if replyTo != nil {
		input.ReplyToAddresses = []string{*replyTo}
	}

	_, err := s.client.SendEmail(ctx, input)
	return err
}

func NewSES(id string) (Provider, error) {
	envId := strings.ToUpper(id)
	region := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_REGION", envId))
	if len(region) == 0 {
		return nil, fmt.Errorf("missing region for ses provider %q", id)
	}
	accessKey := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_ACCESS_KEY", envId))
	if len(accessKey) == 0 {
		return nil, fmt.Errorf("missing access key for ses provider %q", id)
	}
	secretKey := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_SECRET_KEY", envId))
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("missing secret key for ses provider %q", id)
	}
	sessionToken := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_SESSION_TOKEN", envId))

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil, err
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	// Create and verify the client
	client := sesv2.NewFromConfig(cfg)
	if _, err := client.ListEmailIdentities(context.Background(), nil); err != nil {
		return nil, err
	}

	// Get the sending limits
	account, err := client.GetAccount(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	perSecond := int(math.Round(account.SendQuota.MaxSendRate))

	return &SES{
		client: client,
		rl:     ratelimit.New(perSecond, ratelimit.Per(time.Second), ratelimit.WithSlack(10)),
	}, nil
}
