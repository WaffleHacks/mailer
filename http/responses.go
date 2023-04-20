package http

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func success(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func failure(w http.ResponseWriter, code int, reason string) {
	w.WriteHeader(code)
	if _, err := fmt.Fprintf(w, `{"error":"%s"}`, reason); err != nil {
		zap.L().Named("http").Error(
			"failed to send error response",
			zap.Int("response.code", code),
			zap.String("response.reason", reason),
			zap.Error(err),
		)
	}
}

func deserializationFailure(w http.ResponseWriter) {
	failure(w, http.StatusUnprocessableEntity, "failed to deserialize request body")
}
