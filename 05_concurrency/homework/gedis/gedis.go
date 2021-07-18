package gedis

import (
	"sync"
	"time"
)

const cleanupInterval = 5 * time.Second

// Element struct contains value (interface{}) and its TTL
type element struct {
	value       interface{}
	validBefore time.Time
	ttl         time.Duration
}

type Gedis struct {
	data map[string]*element
	sync.Mutex
}

func (g *Gedis) Set(key string, value interface{}, ttl time.Duration) {
	// Check if key already exists
	if e, ok := g.data[key]; ok {
		// ...and update it
		g.Lock()
		defer g.Unlock()
		e.value = value
		e.ttl = ttl
		e.validBefore.Add(ttl)
		return
	}
	// or create
	g.Lock()
	defer g.Unlock()
	e := new(element)
	e.value = value
	e.ttl = ttl
	e.validBefore = time.Now().Add(ttl)
	g.data[key] = e
}

func (g *Gedis) Get(key string) (interface{}, bool) {
	e, found := g.data[key]
	if !found {
		return nil, false
	}
	g.Lock()
	defer g.Unlock()
	e.validBefore = time.Now().Add(e.ttl)
	return e.value, true
}

func (g *Gedis) Delete(key string) {
	_, ok := g.data[key]
	if ok {
		g.Lock()
		defer g.Unlock()
		delete(g.data, key)
	}
}

func cleanup(g *Gedis) {
	for {
		time.Sleep(cleanupInterval)
		for k, v := range g.data {
			if time.Now().After(v.validBefore) {
				g.Lock()
				delete(g.data, k)
				g.Unlock()
			}
		}
	}

}

func NewGedis() *Gedis {
	data := make(map[string]*element)
	g := new(Gedis)
	g.data = data
	go cleanup(g)
	return g
}
