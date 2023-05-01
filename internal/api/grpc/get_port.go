package grpc

import (
	"context"
	"errors"
	"github.com/d6o/portsservice/internal/domain/model"
	"github.com/d6o/portsservice/internal/usecases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	GetPort struct {
		UnimplementedPortServiceServer
		portFinder portFinder
	}

	portFinder interface {
		FindPort(ctx context.Context, portKey string) (*model.Port, error)
	}
)

func NewGetPort(portFinder portFinder) *GetPort {
	return &GetPort{portFinder: portFinder}
}

func (gp *GetPort) GetPort(ctx context.Context, in *GetPortRequest) (*GetPortResponse, error) {
	if in.PortKey == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	port, err := gp.portFinder.FindPort(ctx, in.PortKey)
	if errors.Is(err, usecases.ErrPortDoesNotExist) {
		return nil, status.Errorf(codes.NotFound, "Port not found in the storage")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while trying to find port")
	}

	return portToDTO(port), nil
}

func portToDTO(port *model.Port) *GetPortResponse {
	dto := &GetPortResponse{
		Name:     port.Name,
		City:     port.City,
		Country:  port.Country,
		Alias:    port.Alias,
		Regions:  port.Regions,
		Province: port.Province,
		Timezone: port.Timezone,
		Unlocs:   port.Unlocs,
		Code:     port.Code,
	}

	if len(port.Coordinates) >= 1 {
		dto.Lat = port.Coordinates[0]
	}

	if len(port.Coordinates) >= 2 {
		dto.Lon = port.Coordinates[1]
	}

	return dto
}
