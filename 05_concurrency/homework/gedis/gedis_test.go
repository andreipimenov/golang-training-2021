package gedis

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {

	var cache = NewGedis()

	key := "common_key"
	val := 1

	cache.Set(key, val, 5*time.Second)

	// getting correct values
	v, found := cache.Get(key)
	if !found {
		t.Errorf("Value should be found: %v", key)
	}
	if v != val {
		t.Errorf("Wrong value for %v, get %v, expected %v", key, v, val)
	}
}

func TestExpiredKey(t *testing.T) {
	var cache = NewGedis()
	// lets test if expiry time is increasing
	key := "i_will_expire_soon"
	val := 1
	cache.Set(key, val, 3*time.Second)
	time.Sleep(5 * time.Second)
	v, found := cache.Get(key)
	if v != nil || found != false {
		t.Errorf("Key should NOT be found: key: %v, get: %v", key, val)
	}
}

func TestIncreaseExpiry(t *testing.T) {
	var cache = NewGedis()
	// lets test if expiry time is increasing
	key := "we_will_get_it_from_goroutine"
	val := 1
	cache.Set(key, val, 3*time.Second)
	go func() {
		for {
			cache.Get(key)
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(5 * time.Second)
	v, found := cache.Get(key)
	if v != val || found == false {
		t.Errorf("Key should be found: key: %v, get: %v", key, val)
	}
}

func TestDelete(t *testing.T) {
	var cache = NewGedis()
	key := "to_be_deleted_right_after_set"
	val := 1
	cache.Set(key, val, 10*time.Second)
	cache.Delete(key)
	v, found := cache.Get(key)
	if found {
		t.Errorf("Key should NOT be found: key: %v, get: %v", key, v)
	}
}

func TestWrongKey(t *testing.T) {
	var cache = NewGedis()
	// gtting wrond value
	key := "I'm absent"
	value, found := cache.Get(key)
	if value != nil || found != false {
		t.Errorf("Key should not be found: key: %v, get: %v", key, value)
	}
}
