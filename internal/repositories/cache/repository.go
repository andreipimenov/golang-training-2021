package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

type Repository struct {
	data sync.Map
}

var _ service.CacheRepository = (*Repository)(nil)

func New() *Repository {
	return &Repository{data: sync.Map{}}
}

func formKey(ticker string, date time.Time) string {
	return fmt.Sprintf("%s_%s", ticker, date.Format("2006-01-02"))
}

func (r *Repository) Load(ticker string, date time.Time) (model.Price, bool) {
	value, ok := r.data.Load(formKey(ticker, date))
	if !ok {
		return model.Price{}, false
	}
	p, ok := value.(model.Price)
	return p, ok
}

func (r *Repository) Store(ticker string, date time.Time, price model.Price) {
	r.data.Store(formKey(ticker, date), price)
}
