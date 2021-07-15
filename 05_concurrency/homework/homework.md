# Homework 03

Implement in-memory cache.
It must be safe for concurrent usage.

Cache interface:
```
type Cache interface {
    Set(key string, value interface{}, ttl time.Duration)
    Get(key string) (interface{}, bool)
    Delete(key string)
}
```

TTL means duration while key is valid. Invalidation should happen automatically.
After reading a key its TTL should be increased up to current time + TTL.
