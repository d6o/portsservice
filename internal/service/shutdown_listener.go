package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/d6o/portsservice/internal/infrastructure/appcontext"
	"go.uber.org/zap"
)

type (
	ShutdownListener struct {
		cancel     context.CancelFunc
		closeables []closeable
		done       chan struct{}
	}
	closeable func(ctx context.Context) error
)

func NewShutdownListener(cancel context.CancelFunc) *ShutdownListener {
	return &ShutdownListener{
		cancel:     cancel,
		closeables: []closeable{},
		done:       make(chan struct{}, 1),
	}
}

func (sl *ShutdownListener) Listen(ctx context.Context) {
	terminationSignal := make(chan os.Signal, 1)

	terminationSignals := []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGUSR2,
	}

	signal.Notify(terminationSignal, terminationSignals...)

	receivedSignal := <-terminationSignal

	logger := appcontext.Logger(ctx).With(zap.Any("signal", receivedSignal))
	logger.Warn("Shutdown signal received, closing context")

	sl.cancel()

	logger.Warn("calling closeables")

	for _, c := range sl.closeables {
		if err := c(ctx); err != nil {
			logger.With(zap.Error(err)).Warn("Error when calling closeable")
		}
	}

	sl.done <- struct{}{}
}

func (sl *ShutdownListener) AddClosable(c closeable) {
	sl.closeables = append(sl.closeables, c)
}

func (sl *ShutdownListener) Done() chan struct{} {
	return sl.done
}
