package main

import (
	"testing"
	"time"
)

const (
	testKey string = "test_key"
	testValue string = "test_value"
	testKeyNotPresent = "test_not_present"
)

var cache = New(time.Second, 4*time.Second)

func TestMemCache_Get(t *testing.T) {
	cache.Set(testKey,testValue, 2*time.Minute)

	value, found := cache.Get(testKey)

	if value != testValue {
		t.Error("Error: ", "The received value is not what it should be:", value, testValue)
	}

	if found != true {
		t.Error("Error: ", "Could not get cache")
	}

	//Try to get a value for key not in map
	value, found = cache.Get(testKeyNotPresent)

	if value != nil || found != false {
		t.Error("Error: ", "Value does not exist and must be empty", value)
	}
}

func TestMemCache_Delete(t *testing.T) {
	cache.Set(testKey, testValue, 1*time.Minute)

	cache.Delete(testKey)

	value, found := cache.Get(testKey)

	if found {
		t.Error("Error: ", "Should not be found because it was deleted")
	}

	if value != nil {
		t.Error("Error: ", "Value is not nil:", value)
	}

}