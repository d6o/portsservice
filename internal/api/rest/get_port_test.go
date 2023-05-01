package rest_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/d6o/portsservice/internal/api/rest"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/golang/mock/gomock"
)

func TestGetPort_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	portRetriever := model.NewMockRetriever(ctrl)
	responder := rest.NewMockresponder(ctrl)

	w := httptest.NewRecorder()

	port := &model.Port{
		Name: "Test",
	}

	responder.EXPECT().WriteOK(gomock.Any(), w, port)

	getPort := rest.NewGetPort(portRetriever, responder)

	portRetriever.EXPECT().Get(gomock.Any(), "1").Return(port, true, nil)

	req := httptest.NewRequest(http.MethodGet, "/ports/1", bytes.NewBuffer(nil))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("portKey", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	getPort.ServeHTTP(w, req)
}

func TestGetPort_ServeHTTP_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	portRetriever := model.NewMockRetriever(ctrl)
	responder := rest.NewMockresponder(ctrl)

	w := httptest.NewRecorder()

	responder.EXPECT().WriteNotFound(gomock.Any(), w)

	getPort := rest.NewGetPort(portRetriever, responder)

	portRetriever.EXPECT().Get(gomock.Any(), "not_found").Return(nil, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/ports/not_found", bytes.NewBuffer(nil))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("portKey", "not_found")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	getPort.ServeHTTP(w, req)
}

func TestGetPort_ServeHTTP_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	portRetriever := model.NewMockRetriever(ctrl)
	responder := rest.NewMockresponder(ctrl)
	w := httptest.NewRecorder()

	getPort := rest.NewGetPort(portRetriever, responder)

	portRetriever.EXPECT().Get(gomock.Any(), "1").Return(nil, false, errors.New("test error"))
	responder.EXPECT().WriteInternalServerError(gomock.Any(), w, "can't retrieve port")

	req := httptest.NewRequest(http.MethodGet, "/ports/1", bytes.NewBuffer(nil))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("portKey", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	getPort.ServeHTTP(w, req)
}
