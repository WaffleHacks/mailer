package providers

import (
	"context"
)

// Provider is an interface to an email provider
type Provider interface {
	// Send dispatches an email message to the provider for sending to one recipient
	Send(ctx context.Context, to string, from string, subject string, body string, htmlBody, replyTo *string) error
}

// BatchedProvider is an extension to the Provider stating that the provider natively supports batching
type BatchedProvider interface {
	// SendBatch dispatches an email message to the provider for sending to many recipients
	SendBatch(ctx context.Context, to []string, from string, subject string, body string, htmlBody, replyTo *string) error
}
