package api

import (
	"context"
	"net/http"

	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/d6o/portsservice/internal/infrastructure/appcontext"
	"github.com/go-chi/chi/v5"
)

//go:generate mockgen -source get_port.go -destination mock_get_port.go -package api

type (
	GetPort struct {
		portRetriever model.Retriever
		responder     responder
	}

	responder interface {
		WriteInternalServerError(ctx context.Context, w http.ResponseWriter, msg string)
		WriteBadRequest(ctx context.Context, w http.ResponseWriter, msg string)
		WriteNotFound(ctx context.Context, w http.ResponseWriter)
		WriteOK(ctx context.Context, w http.ResponseWriter, result any)
	}

	GetPortResponse struct {
		Name        string    `json:"name"`
		City        string    `json:"city"`
		Country     string    `json:"country"`
		Alias       []string  `json:"alias"`
		Regions     []string  `json:"regions"`
		Coordinates []float64 `json:"coordinates"`
		Province    string    `json:"province"`
		Timezone    string    `json:"timezone"`
		Unlocs      []string  `json:"unlocs"`
		Code        string    `json:"code"`
	}
)

func NewGetPortResponse(port *model.Port) *GetPortResponse {
	return &GetPortResponse{
		Name:        port.Name,
		City:        port.City,
		Country:     port.Country,
		Alias:       port.Alias,
		Regions:     port.Regions,
		Coordinates: port.Coordinates,
		Province:    port.Province,
		Timezone:    port.Timezone,
		Unlocs:      port.Unlocs,
		Code:        port.Code,
	}
}

func NewGetPort(portRetriever model.Retriever, responder responder) *GetPort {
	return &GetPort{portRetriever: portRetriever, responder: responder}
}

func (gp GetPort) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.Logger(r.Context())
	logger.Info("Request received")

	portKey := chi.URLParam(r, "portKey")
	if len(portKey) == 0 {
		gp.responder.WriteBadRequest(r.Context(), w, "invalid port key")
		return
	}

	port, exist, err := gp.portRetriever.Get(r.Context(), portKey)
	if err != nil {
		gp.responder.WriteInternalServerError(r.Context(), w, "can't retrieve port")
		return
	}

	if !exist {
		gp.responder.WriteNotFound(r.Context(), w)
		return
	}

	dto := NewGetPortResponse(port)

	gp.responder.WriteOK(r.Context(), w, dto)

	logger.Info("Request handled")
}
