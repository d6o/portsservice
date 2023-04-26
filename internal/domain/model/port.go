package model

import "context"

//go:generate mockgen -source port.go -destination mock_port.go -package model

type (
	Port struct {
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
