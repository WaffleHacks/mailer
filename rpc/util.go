package rpc

import (
	"net/mail"

	"github.com/WaffleHacks/mailer/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isEmailValid(address string, l *logging.Logger, message string) *status.Status {
	if _, err := mail.ParseAddress(address); err != nil {
		l.Warn(message, zap.String("email", address))
		return status.New(codes.InvalidArgument, message)
	}

	return nil
}
