package providers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/getsentry/sentry-go"

	"github.com/WaffleHacks/mailer/logging"
)

func newContent(data string) *types.Content {
	return &types.Content{
		Data:    &data,
		Charset: aws.String("UTF-8"),
	}
}

type SES struct {
	client *sesv2.Client
}

func (s *SES) Send(ctx context.Context, _ *logging.Logger, to, from, subject, body string, htmlBody, replyTo *string) error {
	span := sentry.TransactionFromContext(ctx).StartChild("send")
	defer span.Finish()

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

	_, err := s.client.SendEmail(span.Context(), input)
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

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil, err
	}

	// Create and verify the client
	client := sesv2.NewFromConfig(cfg)
	if _, err := client.ListEmailIdentities(context.Background(), nil); err != nil {
		return nil, err
	}

	return &SES{
		client: client,
	}, nil
}
