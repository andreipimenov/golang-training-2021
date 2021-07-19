package memorycache

import (
	"sync"
	"time"
)

const defaultTtl = 5 * time.Second
const cleanInterval = time.Second

type Item struct {
	Value      interface{}
	created    time.Time
	expiration int64
}

type MemoryCache struct {
	sync.RWMutex
	cash map[string]Item
}

func (m *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	m.Lock()
	defer m.Unlock()

	if ttl <= 0 {
		ttl = defaultTtl
	}

	created := time.Now()

	m.cash[key] = Item{
		value,
		created,
		created.Add(ttl).Unix(),
	}
}

func (m *MemoryCache) Get(key string) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()

	item, isHas := m.cash[key]

	if isHas && time.Now().Unix() <= item.expiration {
		return item.Value, true
	}

	return nil, false
}

func (m *MemoryCache) Delete(key string) {
	m.Lock()
	defer m.Unlock()

	if _, isHas := m.cash[key]; isHas {
		delete(m.cash, key)
	}
}

func InitializeMemoryCache() *MemoryCache {
	cash := &MemoryCache{cash: make(map[string]Item)}
	go cash.startGC()
	return cash
}

func (m *MemoryCache) startGC() {
	for {
		<-time.After(cleanInterval)

		if m.cash == nil {
			return
		}

		if keys := m.getExpiredKeys(); len(keys) != 0 {
			m.clearItems(keys)
		}
	}
}

func (m *MemoryCache) getExpiredKeys() (keys []string) {
	m.RLock()
	defer m.RUnlock()

	for key, item := range m.cash {
		if time.Now().Unix() > item.expiration {
			keys = append(keys, key)
		}
	}

	return
}

func (m *MemoryCache) clearItems(keys []string) {
	m.Lock()
	defer m.Unlock()
	for _, k := range keys {
		delete(m.cash, k)
		//fmt.Printf("Clear value for key '%s'\n", k)
	}
}

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}
