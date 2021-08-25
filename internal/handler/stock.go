package handler

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/pb"
)

const (
	StockPath = "/price/{ticker}/{date}"
)

type Stock struct {
	pb.UnimplementedStockServer
	logger  *zerolog.Logger
	service StockService
}

type StockService interface {
	GetPrice(string, time.Time) (*model.Price, error)
}

func NewStock(logger *zerolog.Logger, srv StockService) *Stock {
	return &Stock{
		logger:  logger,
		service: srv,
	}
}

func (s *Stock) GetPrice(ctx context.Context, in *pb.GetPriceRequest) (*pb.GetPriceResponse, error) {
	d, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		s.logger.Error().Err(err).Msg("Invalid incoming date parameter")
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	price, err := s.service.GetPrice(in.Ticker, d)
	if err != nil {
		s.logger.Error().Err(err).Msg("GetPrice method error")
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &pb.GetPriceResponse{
		Open:  price.Open,
		High:  price.High,
		Low:   price.Low,
		Close: price.Close,
	}, nil
}
