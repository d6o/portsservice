package model

import "context"

//go:generate mockgen -source port.go -destination mock_port.go -package model

type (
	Port struct {
		Name        string
		City        string
		Country     string
		Alias       []string
		Regions     []string
		Coordinates []float64
		Province    string
		Timezone    string
		Unlocs      []string
		Code        string
	}

	Storage interface {
		Saver
		Retriever
	}

	Saver interface {
		Save(ctx context.Context, key string, value *Port) error
	}

	Retriever interface {
		Get(ctx context.Context, key string) (*Port, bool, error)
	}
)
