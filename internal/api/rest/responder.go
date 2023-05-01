// Package api is the package to hold all presentation layer in the HTTP API. Is the entrypoint for requests.
package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/d6o/portsservice/internal/infrastructure/appcontext"
	"go.uber.org/zap"
)

type (
	Responder struct{}

	errorResponse struct {
		Message string `json:"Message"`
	}
)

func NewResponder() *Responder {
	return &Responder{}
}

func (r Responder) WriteBadRequest(ctx context.Context, w http.ResponseWriter, msg string) {
	logger := appcontext.Logger(ctx).With(zap.String("msg", msg))
	logger.Warn("Writing a Bad Request response")

	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(errorResponse{msg}); err != nil {
		logger.With(zap.Error(err)).Warn("Can't write response")
	}
}

func (r Responder) WriteOK(ctx context.Context, w http.ResponseWriter, result any) {
	logger := appcontext.Logger(ctx).With(zap.Any("result", result))
	logger.Info("Writing an OK response")

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.With(zap.Error(err)).Warn("Can't write response")
	}
}

func (r Responder) WriteOKPlain(ctx context.Context, w http.ResponseWriter, result []byte) {
	logger := appcontext.Logger(ctx).With(zap.Any("result", result))
	logger.Info("Writing an OK response")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(result); err != nil {
		logger.With(zap.Error(err)).Warn("Can't write response")
	}
}

func (r Responder) WriteInternalServerError(ctx context.Context, w http.ResponseWriter, msg string) {
	logger := appcontext.Logger(ctx).With(zap.String("msg", msg))
	logger.Warn("Writing a Internal Server Error response")

	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(errorResponse{msg}); err != nil {
		logger.With(zap.Error(err)).Warn("Can't write response")
	}
}

func (r Responder) WriteNotFound(ctx context.Context, w http.ResponseWriter) {
	appcontext.Logger(ctx).Warn("Writing a Not Found response")
	w.WriteHeader(http.StatusNotFound)
}
