package interfaces

import "time"

type Key = string

type Cache interface {
	Set(key Key, value interface{}, ttl time.Duration)
	Get(key Key) (interface{}, bool)
	Delete(key Key)

	// Cancels any scheduled sweeping.
	// Some cache implementations may require to call Stop to avoid memory leaks.
	// After calling Stop, further usage of the cache is invalid.
	Stop()
}
