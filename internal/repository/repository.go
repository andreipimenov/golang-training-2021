package repository

import (
	"sync"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type Repository struct {
	data sync.Map
}

func New() *Repository {
	return &Repository{data: sync.Map{}}
}

func (r *Repository) Load(key string) (model.Ticker, bool) {
	value, ok := r.data.Load(key)
	if !ok {
		return model.Ticker{}, false
	}
	p, ok := value.(model.Ticker)
	return p, ok
}

func (r *Repository) Store(key string, value model.Ticker) {
	r.data.Store(key, value)
}
