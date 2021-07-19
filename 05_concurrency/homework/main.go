package main

import (
	"sync"
	"time"
)

type cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type item struct {
	Value      interface{}
	Expiration int64
	TTL        time.Duration
}

type memoryCache struct {
	sync.Mutex
	defaultExpirationTime time.Duration
	cleanUpInterval       time.Duration
	items                 map[string]item
}

func createCache(defaultExpirationTime, cleanUpInterval time.Duration) *memoryCache {
	items := make(map[string]item)

	cache := memoryCache{
		items:                 items,
		defaultExpirationTime: defaultExpirationTime,
		cleanUpInterval:       cleanUpInterval,
	}

	if defaultExpirationTime <= 0 {
		defaultExpirationTime = time.Second
	}

	if cleanUpInterval > 0 {
		cache.startCleanUp()
	}

	return &cache
}

func (m *memoryCache) set(key string, value interface{}, ttl time.Duration) {
	var expiration int64

	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	} else {
		ttl = m.defaultExpirationTime
	}

	m.Lock()
	defer m.Unlock()

	m.items[key] = item{
		Value:      value,
		Expiration: expiration,
		TTL:        ttl,
	}
}

func (m *memoryCache) get(key string) (interface{}, bool) {
	m.Lock()
	defer m.Unlock()

	mValue, found := m.items[key]

	if !found {
		return nil, false
	}

	if mValue.Expiration > 0 {
		if time.Now().UnixNano() > mValue.Expiration {
			return nil, false
		}
	}

	mValue.Expiration = time.Now().Add(mValue.TTL).UnixNano()

	return mValue.Value, true

}

func (m *memoryCache) getExpiredKeys() (keys []string) {
	m.Lock()
	defer m.Unlock()

	for key, item := range m.items {
		if time.Now().UnixNano() > item.Expiration && item.Expiration > 0 {
			keys = append(keys, key)
		}
	}

	return
}

func (m *memoryCache) delete(key string) {
	m.Lock()
	defer m.Unlock()

	delete(m.items, key)
}

func (m *memoryCache) cleanUp() {
	for {
		<-time.After(m.cleanUpInterval)

		if m.items == nil {
			return
		}

		if keys := m.getExpiredKeys(); len(keys) != 0 {
			m.Lock()
			for _, key := range keys {
				delete(m.items, key)
			}
			m.Unlock()
		}

	}
}

func (m *memoryCache) startCleanUp() {
	go m.cleanUp()
}

