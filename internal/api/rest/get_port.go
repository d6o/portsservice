package rest

import (
	"context"
	"errors"
	"github.com/d6o/portsservice/internal/usecases"
	"net/http"

	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/d6o/portsservice/internal/infrastructure/appcontext"
	"github.com/go-chi/chi/v5"
)

//go:generate mockgen -source get_port.go -destination mock_get_port.go -package api

type (
	GetPort struct {
		portFinder portFinder
		responder  responder
	}

	responder interface {
		WriteInternalServerError(ctx context.Context, w http.ResponseWriter, msg string)
		WriteBadRequest(ctx context.Context, w http.ResponseWriter, msg string)
		WriteNotFound(ctx context.Context, w http.ResponseWriter)
		WriteOK(ctx context.Context, w http.ResponseWriter, result any)
	}

	portFinder interface {
		FindPort(ctx context.Context, portKey string) (*model.Port, error)
	}
)

func NewGetPort(portFinder portFinder, responder responder) *GetPort {
	return &GetPort{portFinder: portFinder, responder: responder}
}

func (gp GetPort) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.Logger(r.Context())
	logger.Info("Request received")

	portKey := chi.URLParam(r, "portKey")
	if len(portKey) == 0 {
		gp.responder.WriteBadRequest(r.Context(), w, "invalid port key")
		return
	}

	port, err := gp.portFinder.FindPort(r.Context(), portKey)
	if errors.Is(err, usecases.ErrPortDoesNotExist) {
		gp.responder.WriteNotFound(r.Context(), w)
		return
	}
	if err != nil {
		gp.responder.WriteInternalServerError(r.Context(), w, "can't retrieve port")
		return
	}

	gp.responder.WriteOK(r.Context(), w, port)

	logger.Info("Request handled")
}
