package checkall

import (
	"memcache/interfaces"
	"sync"
	"time"
)

type cache struct {
	mu      sync.Mutex
	records map[interfaces.Key]*record
	stop    chan struct{}
}

type record struct {
	value interface{}
	ttl   time.Duration
	birth time.Time
}

func New(sweepInterval time.Duration) interfaces.Cache {
	c := &cache{
		records: make(map[interfaces.Key]*record),
		stop:    make(chan struct{}),
	}
	go func() {
		ticker := time.NewTicker(sweepInterval)
		defer ticker.Stop()
		for {
			select {
			case <-c.stop:
				return
			case now := <-ticker.C:
				c.mu.Lock()
				for k, r := range c.records {
					death := r.birth.Add(r.ttl)
					if now.After(death) {
						delete(c.records, k)
					}
				}
				c.mu.Unlock()
			}
		}
	}()
	return c
}

func (c *cache) Stop() {
	close(c.stop)
}

func (c *cache) Set(key interfaces.Key, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.records[key] = &record{value, ttl, time.Now()}
}

func (c *cache) Get(key interfaces.Key) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record, ok := c.records[key]
	if ok {
		record.birth = time.Now()
		value = record.value
	}
	return
}

func (c *cache) Delete(key interfaces.Key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.records, key)
}
