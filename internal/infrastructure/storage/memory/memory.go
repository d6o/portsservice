package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/d6o/portsservice/internal/domain/model"
)

//go:generate mockgen -source memory.go -destination mock_memory.go -package memory

type (
	InMemory struct {
		data           sync.Map
		contextChecker contextChecker
	}

	contextChecker interface {
		CheckContext(ctx context.Context) error
	}
)

func NewInMemory(contextChecker contextChecker) *InMemory {
	return &InMemory{
		data:           sync.Map{},
		contextChecker: contextChecker,
	}
}

func (db *InMemory) Save(ctx context.Context, key string, port *model.Port) error {
	if err := db.contextChecker.CheckContext(ctx); err != nil {
		return err
	}

	db.data.Store(key, port)

	return nil
}

func (db *InMemory) Get(ctx context.Context, key string) (*model.Port, bool, error) {
	if err := db.contextChecker.CheckContext(ctx); err != nil {
		return nil, false, err
	}

	val, ok := db.data.Load(key)
	if !ok {
		return nil, ok, nil
	}

	port, ok := val.(*model.Port)
	if !ok {
		return nil, false, errors.New("invalid object")
	}

	return port, true, nil
}
