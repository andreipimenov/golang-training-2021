package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type Self struct{}

func (s *Self) Get(key string) (interface{}, bool) {
	return nil, false
}

func (s *Self) Set(key string, value interface{}, ttl time.Duration) {}

func (s *Self) Delete(key string) {}

type Element struct {
	value interface{}
	ttl   time.Duration
	ttp   time.Time
}

type CacheMemory struct {
	sync.Mutex
	object map[string]Element
}

func NewCacheMemory() *CacheMemory {
	cacheNew := &CacheMemory{object: make(map[string]Element)}
	cacheNew.openAssertion()
	return cacheNew
}
func (m *CacheMemory) Assertion() {
	for {
		<-time.After(5 * time.Second)
		element := m.oldCache()
		if len(element) > 0 {
			m.cleanCache(element)
		}
	}
}
func (m *CacheMemory) openAssertion() {
	go m.Assertion()
}

func (m *CacheMemory) cleanCache(element []string) {
	for _, e := range element {
		m.Delete(e)
	}
}

func (m *CacheMemory) oldCache() (oldElem []string) {
	for element, j := range m.object {
		if time.Now().UnixNano() > int64(j.ttl.Seconds()) {
			oldElem = append(oldElem, element)
		}
	}
	return oldElem
}
func (m *CacheMemory) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.object, key)
}

func (m *CacheMemory) Get(key string) (interface{}, bool) {
	m.Lock()
	defer m.Unlock()
	element, ok := m.object[key]
	if !ok {
		return nil, false
	}
	m.object[key] = Element{
		element.value,
		m.object[key].ttl,
		time.Now(),
	}
	return element.value, true
}

func (m *CacheMemory) Set(key string, value interface{}, ttl time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.object[key] = Element{
		value,
		ttl,
		time.Now(),
	}
}

func main() {
	ourCache := NewCacheMemory()
	ourCache.Set("firstKey", "firstValue", 5*time.Second)
	time.Sleep(4 * time.Second)
	fmt.Println(ourCache.Get("firstKey"))
	ourCache.Set("secondKey", "secondValue", 3*time.Second)
	time.Sleep(3 * time.Second)
	fmt.Println(ourCache.Get("secondKey"))
	ourCache.Set("thirdKey", "thirdValue", 1*time.Second)
	fmt.Println(ourCache.Get("thirdKey"))
	ourCache.Delete("thirdKey")
	time.Sleep(10 * time.Second)
	fmt.Println("firstKey", "secondKey")

}
