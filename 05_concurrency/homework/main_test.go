package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type Value struct {
	value interface{}
	ttl   time.Duration
	exp   time.Time
}

func (v Value) IsExpired() bool {
	return v.exp.Sub(time.Now()) < 0
}

func (v Value) String() string {
	return "> value: " + fmt.Sprintf("%v", v.value) + "; " +
		"ttl: " + v.ttl.String() + "\n"
}

type ActionType int32

func (t ActionType) Lookup() string {
	switch t {
	case Set:
		return "Set"
	case Delete:
		return "Delete"
	case Read:
		return "Read"
	}
	return "unknown"
}

const (
	Set ActionType = iota
	Delete
	Read
)

type Action struct {
	key         string
	actionValue Value
	actionType  ActionType
	resp        chan ReadResp
}

func (a Action) String() string {
	return "> action:\n" +
		"key: " + a.key + "\n" +
		"actionValue: " + a.actionValue.String() +
		"actionType: " + a.actionType.Lookup() + "\n"
}

type ReadResp struct {
	ok        bool
	respValue Value
}

type InMemoryCache struct {
	values  map[string]Value
	actions chan Action
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.actions <- Action{
		key: key,
		actionValue: Value{
			value: value,
			ttl:   ttl,
			exp:   time.Now().Add(ttl),
		},
		actionType: Set,
	}
}

func (c *InMemoryCache) Delete(key string) {
	c.actions <- Action{
		key:        key,
		actionType: Delete,
	}
}

func (c *InMemoryCache) GetResp(key string) ReadResp {
	ch := make(chan ReadResp)
	c.actions <- Action{
		key:        key,
		resp:       ch,
		actionType: Read,
	}
	for {
		if r, ok := <-ch; ok {
			return r
		}
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	r := c.GetResp(key)
	if r.respValue.IsExpired() {
		c.Delete(key)
		return nil, false
	}
	c.Set(key, r.respValue.value, r.respValue.ttl)
	return r.respValue.value, r.ok
}

func (c *InMemoryCache) initActions() *InMemoryCache {
	go func() {
		for {
			c.proceedAction(<-c.actions)
		}
	}()
	return c
}

func (c *InMemoryCache) proceedAction(a Action) {
	//fmt.Println(a)
	switch a.actionType {
	case Set:
		c.values[a.key] = a.actionValue
		c.ScheduleCheckTtl(a.key, a.actionValue.ttl)
	case Delete:
		delete(c.values, a.key)
	case Read:
		value, ok := c.values[a.key]
		a.resp <- ReadResp{
			ok:        ok,
			respValue: value,
		}
	}
}

func (c *InMemoryCache) ScheduleCheckTtl(key string, ttl time.Duration) {
	go func() {
		time.Sleep(ttl)
		c.CheckTtl(key)
	}()
}

func (c *InMemoryCache) CheckTtl(key string) {
	r := c.GetResp(key)
	if r.respValue.IsExpired() {
		c.Delete(key)
	}
}

func NewInMemoryCache() *InMemoryCache {
	c := InMemoryCache{
		values:  map[string]Value{},
		actions: make(chan Action),
	}
	return c.initActions()
}

func Test_WriteRead(t *testing.T) {
	run := []int{
		1,
	}
	cache := NewInMemoryCache()
	for _, test := range run {
		cache.Set(strconv.Itoa(test), test, time.Minute)
		time.Sleep(time.Second)
		r, _ := cache.Get(strconv.Itoa(test))
		if r != test {
			t.Errorf("Expected %v but got %v", test, r)
			return
		}
	}
}

func Test_TtlExpire(t *testing.T) {
	run := []int{
		1,
	}
	cache := NewInMemoryCache()
	for _, test := range run {
		cache.Set(strconv.Itoa(test), test, time.Second)
		time.Sleep(2 * time.Second)
		_, ok := cache.Get(strconv.Itoa(test))
		if ok != false {
			t.Errorf("Expected %v but got %v", false, ok)
			return
		}
	}
}

func Test_TtlNotExpire(t *testing.T) {
	run := []int{
		1,
	}
	cache := NewInMemoryCache()
	for _, test := range run {
		cache.Set(strconv.Itoa(test), test, 2*time.Second)
		time.Sleep(time.Second)
		_, ok := cache.Get(strconv.Itoa(test))
		if ok != true {
			t.Errorf("Expected %v but got %v", true, ok)
			return
		}
	}
}

func Test_TtlUpdated(t *testing.T) {
	run := []int{
		1,
	}
	cache := NewInMemoryCache()
	for _, test := range run {
		cache.Set(strconv.Itoa(test), test, 4*time.Second) //4 seconds ttl
		time.Sleep(3 * time.Second)                        //<--
		cache.Get(strconv.Itoa(test))                      //   |-- 5 seconds wait total
		time.Sleep(2 * time.Second)                        //<--
		_, ok := cache.Get(strconv.Itoa(test))
		if ok != true {
			t.Errorf("Expected %v but got %v", true, ok)
			return
		}
	}
}
