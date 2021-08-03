package service

import (
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type CacheRepository interface {
	Load(ticker string, date time.Time) (model.Price, bool)
	Store(ticker string, date time.Time, price model.Price)
}

type SlowRepository interface {
	GetPrice(ticker string, date time.Time) (model.Price, error)
}

type Service struct {
	cacheRepo CacheRepository
	slowRepo  SlowRepository
}

func New(cacheRepo CacheRepository, slowRepo SlowRepository) *Service {
	return &Service{
		cacheRepo: cacheRepo,
		slowRepo:  slowRepo,
	}
}

func (s *Service) GetPrice(ticker string, date time.Time) (model.Price, error) {
	if p, ok := s.cacheRepo.Load(ticker, date); ok {
		return p, nil
	}
	p, err := s.slowRepo.GetPrice(ticker, date)
	if err != nil {
		return model.Price{}, err
	}
	s.cacheRepo.Store(ticker, date, p)
	return p, nil
}
