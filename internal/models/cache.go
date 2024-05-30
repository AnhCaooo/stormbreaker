package models

import (
	"sync"
	"time"
)

type Cache struct {
	data map[string]cacheValue
	lock sync.Mutex
}

type cacheValue struct {
	value      interface{}
	expiration time.Time
}

// a function is initialize a cache instance
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]cacheValue),
	}
}

// a method is used to add new key-value pair to the cache.
// It takes in a key, a value, and a duration representing the expiration time of the value.
// It first acquires a lock on the mutex to ensure thread safety, and then it adds the key-value pair to the map along with the expiration time.
// Finally, it releases the lock.
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	expirationTime := time.Now().Add(expiration)
	c.data[key] = cacheValue{
		value:      value,
		expiration: expirationTime,
	}
}

// a method is used to retrieve a value from the cache by using a key
// It first acquires a lock on the mutex to ensure thread safety.
// Then checks if the cache contains a value for the given key and if that value has not expired.
// If the value is still valid, it returns the value and a boolean value of `true` to indicate that a valid value was found.
// If the value is not valid or the cache does not contain a value for the given key, it returns `nil` and a boolean value of `false`.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.data[key]
	if !ok || time.Now().After(value.expiration) {
		delete(c.data, key)
		return nil, false
	}

	return value.value, true
}
