package service

import (
	"context"
	"errors"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShutdownListener_Listen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownListener := NewShutdownListener(cancel)

	var closeableCalled bool
	shutdownListener.AddClosable(func(ctx context.Context) error {
		closeableCalled = true
		return nil
	})

	go shutdownListener.Listen(ctx)

	time.Sleep(50 * time.Millisecond) // Give some time to start listening
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	select {
	case <-shutdownListener.Done():
	case <-time.After(1 * time.Second):
		t.Error("ShutdownListener did not finish in a reasonable time")
	}

	assert.True(t, closeableCalled, "Closeable was not called")
}

func TestShutdownListener_ErrorInCloseable(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownListener := NewShutdownListener(cancel)

	var wg sync.WaitGroup
	wg.Add(1)
	shutdownListener.AddClosable(func(ctx context.Context) error {
		defer wg.Done()
		return errors.New("test error")
	})

	go shutdownListener.Listen(ctx)

	time.Sleep(50 * time.Millisecond) // Give some time to start listening
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	select {
	case <-shutdownListener.Done():
	case <-time.After(1 * time.Second):
		t.Error("ShutdownListener did not finish in a reasonable time")
	}

	wg.Wait()
}

func TestShutdownListener_MultipleCloseables(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownListener := NewShutdownListener(cancel)

	var closeable1Called, closeable2Called bool
	shutdownListener.AddClosable(func(ctx context.Context) error {
		closeable1Called = true
		return nil
	})
	shutdownListener.AddClosable(func(ctx context.Context) error {
		closeable2Called = true
		return nil
	})

	go shutdownListener.Listen(ctx)

	time.Sleep(50 * time.Millisecond) // Give some time to start listening
	assert.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))

	select {
	case <-shutdownListener.Done():
	case <-time.After(1 * time.Second):
		t.Error("ShutdownListener did not finish in a reasonable time")
	}

	assert.True(t, closeable1Called, "Closeable1 was not called")
	assert.True(t, closeable2Called, "Closeable2 was not called")
}

func TestShutdownListener_AddClosable(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	shutdownListener := NewShutdownListener(cancel)

	shutdownListener.AddClosable(func(ctx context.Context) error {
		return nil
	})

	assert.Len(t, shutdownListener.closeables, 1, "Closeables should have 1 element")
}
