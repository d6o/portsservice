package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/d6o/portsservice/internal/api"
	"github.com/d6o/portsservice/internal/config"
	domainService "github.com/d6o/portsservice/internal/domain/service"
	"github.com/d6o/portsservice/internal/infrastructure/appcontext"
	"github.com/d6o/portsservice/internal/infrastructure/log"
	"github.com/d6o/portsservice/internal/infrastructure/storage"
	"github.com/d6o/portsservice/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	App struct{}
)

func NewApp() *App {
	return &App{}
}

const (
	minArgs = 2
)

func (App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(os.Args) < minArgs {
		return errors.New("filename was not provided")
	}

	logger, err := log.NewZap()
	if err != nil {
		return err
	}

	defer func() {
		_ = logger.Sync()
	}()

	ctx = appcontext.WithLogger(ctx, logger)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger.With(zap.Any("cfg", cfg)).Info("Configurations loaded")

	contextChecker := appcontext.NewChecker()

	db, err := storage.New(ctx, contextChecker, cfg.StorageProvider, cfg.RedisAddress)
	if err != nil {
		return err
	}

	shutdownListenerService := service.NewShutdownListener(cancel)
	parsePortsDomainService := domainService.NewParsePorts(db, contextChecker)

	go shutdownListenerService.Listen(ctx)

	logger.Info("Importing data to database")

	serverErr := make(chan error, 1)

	go func() {
		filename := os.Args[1]
		if err := parsePortsDomainService.Parse(ctx, filename); err != nil {
			logger.With(zap.Error(err)).Error("Can't parse file")
			serverErr <- err
		}

		logger.Info("Database loaded")
	}()

	logger.With(zap.Int("port", cfg.APIPort)).Info("Starting HTTP Server")

	responder := api.NewResponder()
	getPortHandler := api.NewGetPort(db, responder)
	httpRouter := chi.NewRouter()
	httpRouter.Get("/port/{portKey}", getPortHandler.ServeHTTP)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.APIPort),
		Handler: httpRouter,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	shutdownListenerService.AddClosable(httpServer.Shutdown)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-shutdownListenerService.Done():
		logger.Info("Server shutdown completed")

	case err := <-serverErr:
		logger.Error("Server error occurred", zap.Error(err))
		return err
	}

	return nil
}
