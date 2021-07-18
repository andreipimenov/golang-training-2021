package afterfunc

import (
	"memcache/interfaces"
	"sync"
	"time"
)

type cache struct {
	mu      sync.Mutex
	records map[interfaces.Key]*record
}

type record struct {
	value       interface{}
	ttl         time.Duration
	deleteTimer *time.Timer
}

func New() interfaces.Cache {
	return &cache{
		records: make(map[interfaces.Key]*record),
	}
}

func (c *cache) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, record := range c.records {
		record.deleteTimer.Stop()
	}
}

func (c *cache) Set(key interfaces.Key, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	r, ok := c.records[key]
	if !ok {
		r = &record{}
		c.records[key] = r
		r.deleteTimer = time.AfterFunc(ttl, func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			delete(c.records, key)
		})
	} else {
		r.deleteTimer.Reset(ttl)
	}
	r.value, r.ttl = value, ttl
}

func (c *cache) Get(key interfaces.Key) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record, ok := c.records[key]
	if ok {
		record.deleteTimer.Reset(record.ttl)
		value = record.value
	}
	return
}

func (c *cache) Delete(key interfaces.Key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if record, ok := c.records[key]; ok {
		record.deleteTimer.Stop()
		delete(c.records, key)
	}
}
