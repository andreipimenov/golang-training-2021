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

func (r *Repository) Load(key string) (model.Price, bool) {
	value, ok := r.data.Load(key)
	if !ok {
		return model.Price{}, false
	}
	p, ok := value.(model.Price)
	return p, ok
}

func (r *Repository) Store(key string, value model.Price) {
	r.data.Store(key, value)
}
