package api

import (
	_ "embed"
	"net/http"

	"github.com/d6o/portsservice/internal/infrastructure/appcontext"
)

type (
	GetDocs struct {
	}
)

func NewGetDocs() *GetDocs {
	return &GetDocs{}
}

//go:embed resources/redoc-static.html
var docsPage string

func (gd GetDocs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.Logger(r.Context())
	logger.Info("Request received")

	w.Write([]byte(docsPage))

	logger.Info("Request handled")
}
