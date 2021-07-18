package memcache

import (
	"fmt"
	"math/rand"
	"memcache/afterfunc"
	"memcache/checkall"
	"memcache/interfaces"
	"sync"
	"testing"
	"time"
)

func assertPresence(c interfaces.Cache, key interfaces.Key, value interface{}) {
	if v, ok := c.Get(key); !ok || v != value {
		panic(fmt.Sprint(v, ok))
	}
}
func assertAbsence(c interfaces.Cache, key interfaces.Key) {
	if _, ok := c.Get(key); ok {
		panic("")
	}
}

func cacheCreators(sweepInterval time.Duration) map[string]func() interfaces.Cache {
	return map[string]func() interfaces.Cache{
		"afterfunc": afterfunc.New,
		"checkall":  func() interfaces.Cache { return checkall.New(sweepInterval) },
	}
}

func TestBaseLogic(t *testing.T) {
	long := time.Minute
	for implName, newCache := range cacheCreators(100 * time.Millisecond) {
		t.Run(implName, func(t *testing.T) {
			c := newCache()
			defer c.Stop()
			c.Set("a", 9, long)
			assertPresence(c, "a", 9)
			c.Set("a", 3, long)
			assertPresence(c, "a", 3)
			c.Delete("a")
			assertAbsence(c, "a")
		})
	}
}

func TestDeletingByTime(t *testing.T) {
	for implName, newCache := range cacheCreators(10 * time.Millisecond) {
		t.Run(implName, func(t *testing.T) {
			c := newCache()
			defer c.Stop()

			c.Set("a", 1, 50*time.Millisecond)
			c.Set("b", 2, 150*time.Millisecond)
			assertPresence(c, "a", 1)
			assertPresence(c, "b", 2)

			time.Sleep(75 * time.Millisecond)
			assertAbsence(c, "a")
			assertPresence(c, "b", 2)

			time.Sleep(140 * time.Millisecond)
			assertPresence(c, "b", 2)
			c.Set("b", 3, 50*time.Millisecond)
			assertPresence(c, "b", 3)

			time.Sleep(40 * time.Millisecond)
			assertPresence(c, "b", 3)

			time.Sleep(100 * time.Millisecond)
			assertAbsence(c, "b")
		})
	}
}

// For `go test -race`
func TestConcurrentUsage(t *testing.T) {
	for implName, newCache := range cacheCreators(10 * time.Millisecond) {
		t.Run(implName, func(t *testing.T) {
			c := newCache()
			defer c.Stop()
			wg := sync.WaitGroup{}
			goroutines := 1000
			for i := 0; i < goroutines; i++ {
				wg.Add(1)
				i := i
				go func() {
					k := fmt.Sprint(i / (goroutines / 2)) // "0" or "1"
					v, _ := c.Get(k)
					if v == nil {
						v = 0
					}
					c.Set(k, v.(int)+1, time.Minute)
					wg.Done()
				}()
			}
			wg.Wait()
		})
	}
}

func BenchmarkSet(b *testing.B) {
	long := time.Hour
	rand.Seed(time.Now().UnixNano())
	for implName, newCache := range cacheCreators(long) {
		b.Run(implName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := newCache()
				defer c.Stop()
				for j := 0; j < 10_000; j++ {
					c.Set(fmt.Sprint(j), j, long+time.Duration(rand.Int31()))
				}
			}
		})
	}
}
