package cache

import "github.com/AnhCaooo/stormbreaker/internal/models"

var Cache *models.Cache

// initialize a cache instance
func NewCache() {
	Cache = &models.Cache{
		Data: make(map[string]models.CacheValue),
	}
}
