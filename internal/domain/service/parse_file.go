package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/d6o/portsservice/internal/infrastructure/appcontext"
	"go.uber.org/zap"
)

//go:generate mockgen -source parse_file.go -destination mock_parse_file.go -package service

type (
	ParsePorts struct {
		storage        model.Saver
		contextChecker contextChecker
	}

	contextChecker interface {
		CheckContext(ctx context.Context) error
	}

	portDTO struct {
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

func (p portDTO) ToModel() *model.Port {
	return &model.Port{
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Coordinates: p.Coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		Code:        p.Code,
	}
}

func NewParsePorts(storage model.Saver, contextChecker contextChecker) *ParsePorts {
	return &ParsePorts{storage: storage, contextChecker: contextChecker}
}

func (pp ParsePorts) Parse(ctx context.Context, filename string) error {
	logger := appcontext.Logger(ctx).With(zap.String("file", filename))
	logger.Info("Parsing file")

	if err := pp.contextChecker.CheckContext(ctx); err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			appcontext.Logger(ctx).With(zap.String("file", filename), zap.Error(err)).Warn("can't close file")
		}
	}()

	decoder := json.NewDecoder(bufio.NewReader(file))
	if _, err = decoder.Token(); err != nil {
		return err
	}

	// Loop through each key-value pair in the object
	for decoder.More() {
		if err := pp.contextChecker.CheckContext(ctx); err != nil {
			return err
		}

		key, err := decoder.Token()
		if err != nil {
			return err
		}

		// Decode the value into a Port struct
		var port portDTO
		err = decoder.Decode(&port)
		if err != nil {
			return err
		}

		keyValue, ok := key.(string)
		if !ok {
			return errors.New("key is not valid")
		}

		if err := pp.storage.Save(ctx, keyValue, port.ToModel()); err != nil {
			return err
		}

		logger.With(zap.Any("port", port), zap.String("key", keyValue)).Info("Port saved in storage")
	}

	return nil
}
