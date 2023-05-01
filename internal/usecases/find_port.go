package usecases

import (
	"context"
	"errors"

	"github.com/d6o/portsservice/internal/domain/model"
)

type (
	FindPort struct {
		portRetriever model.Retriever
	}
)

func NewFindPort(portRetriever model.Retriever) *FindPort {
	return &FindPort{portRetriever: portRetriever}
}

var (
	ErrPortDoesNotExist = errors.New("port does not exist")
)

func (fp FindPort) FindPort(ctx context.Context, portKey string) (*model.Port, error) {
	port, exist, err := fp.portRetriever.Get(ctx, portKey)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrPortDoesNotExist
	}

	return port, nil
}
