// AnhCao 2024
package cache

import (
	"sync"
	"time"

	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"go.uber.org/zap"
)

type Cache struct {
	Data   map[string]CacheValue
	logger *zap.Logger
	lock   sync.Mutex
}

type CacheValue struct {
	Value      interface{}
	Expiration time.Time
}

// NewCache returns a new Cache instance
func NewCache(logger *zap.Logger) *Cache {
	return &Cache{
		Data:   make(map[string]CacheValue),
		logger: logger,
	}
}

// a method is used to add new key-value pair to the cache.
// It takes in a key, a value, and a duration representing the expiration time of the value.
// It first acquires a lock on the mutex to ensure thread safety, and then it adds the key-value pair to the map along with the expiration time.
// Finally, it releases the lock.
func (c *Cache) SetExpiredAfterTimePeriod(key string, value interface{}, duration time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	now, err := helpers.GetCurrentTimeInUTC()
	if err != nil {
		c.logger.Fatal("failed to get current time", zap.Error(err))
	}
	expirationTime := now.Add(duration)
	c.logger.Debug("time information",
		zap.Time("current-time-in-utc-zone", time.Now().UTC()),
		zap.Time("expired-time-utc", expirationTime),
	)
	c.Data[key] = CacheValue{
		Value:      value,
		Expiration: expirationTime,
	}
}

// a method is used to add new key-value pair to the cache.
// It takes in a key, a value, and a time slot (by hour) representing the expiration time of the value
// It first acquires a lock on the mutex to ensure thread safety, and then it adds the key-value pair to the map along with the expiration time.
// Finally, it releases the lock.
func (c *Cache) SetExpiredAtTime(key string, value interface{}, expiredTime time.Time) {
	c.logger.Debug("set cache to be expired at",
		zap.Time("expired-time-utc", expiredTime),
		zap.Time("current-time-in-utc-zone", time.Now().UTC()),
	)
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Data[key] = CacheValue{
		Value:      value,
		Expiration: expiredTime,
	}
}

// a method is used to retrieve a value from the cache by using a key
// It first acquires a lock on the mutex to ensure thread safety.
// Then checks if the cache contains a value for the given key and if that value has not expired.
// If the value is still valid, it returns the value and a boolean value of `true` to indicate that a valid value was found.
// If the value is not valid (means not yet cached), it returns `nil` and a boolean value of `false`.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.logger.Debug("get cache key: ", zap.String("key", key))
	c.lock.Lock()
	defer c.lock.Unlock()

	value, exists := c.Data[key]
	if !exists {
		c.logger.Debug("cache key was not found from cache")
		return nil, false
	}
	if time.Now().After(value.Expiration) {
		c.logger.Debug("cache was expired",
			zap.Time("expiration-time-in-utc-zone", value.Expiration),
			zap.Time("current-time-in-utc-zone", time.Now()),
		)
		return nil, false
	}
	c.logger.Debug("cache living time.",
		zap.Any("expired-time-in-utc-zone", value.Expiration),
		zap.Time("current-time-in-utc-zone", time.Now().UTC()),
	)
	return value.Value, true
}

// Delete cache based on receiving cache key. If key is not valid, then Delete is no-op
func (c *Cache) Delete(key string) {
	delete(c.Data, key)

}
