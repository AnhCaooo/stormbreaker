package models

import (
	"sync"
	"time"
)

type Cache struct {
	Data map[string]CacheValue
	lock sync.Mutex
}

type CacheValue struct {
	Value      interface{}
	Expiration time.Time
}

// todo: implement to get Finnish hour instead of current location
// a method is used to add new key-value pair to the cache.
// It takes in a key, a value, and a duration representing the expiration time of the value.
// It first acquires a lock on the mutex to ensure thread safety, and then it adds the key-value pair to the map along with the expiration time.
// Finally, it releases the lock.
func (c *Cache) SetExpiredAfterTimePeriod(key string, value interface{}, expiration time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	expirationTime := time.Now().Add(expiration)
	c.Data[key] = CacheValue{
		Value:      value,
		Expiration: expirationTime,
	}
}

// a method is used to add new key-value pair to the cache.
// It takes in a key, a value, and a time slot (by hour) representing the expiration time of the value (expired at specific hour).
// It first acquires a lock on the mutex to ensure thread safety, and then it adds the key-value pair to the map along with the expiration time.
// Finally, it releases the lock.
func (c *Cache) SetExpiredAtHour(key string, value interface{}, hour int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Get current time
	now := time.Now()

	// Get year, month, and day components
	year, month, day := now.Date()

	// Get today's date
	expiredTime := time.Date(year, month, day, hour, 0, 0, 0, now.Location())

	expirationTime := expiredTime
	c.Data[key] = CacheValue{
		Value:      value,
		Expiration: expirationTime,
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

	value, isValid := c.Data[key]
	if !isValid || time.Now().After(value.Expiration) {
		delete(c.Data, key)
		return nil, false
	}

	return value.Value, true
}
