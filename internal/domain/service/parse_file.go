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
)

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
		var port model.Port
		err = decoder.Decode(&port)
		if err != nil {
			return err
		}

		keyValue, ok := key.(string)
		if !ok {
			return errors.New("key is not valid")
		}

		if err := pp.storage.Save(ctx, keyValue, &port); err != nil {
			return err
		}

		logger.With(zap.Any("port", port), zap.String("key", keyValue)).Info("Port saved in storage")
	}

	return nil
}
