package services

import (
	appLog "api/log"
	"sync"
	"time"
)

// QueryCache provides a simple in-memory cache for database query results
type QueryCache struct {
	cache      map[string]cacheEntry
	mutex      sync.RWMutex
	defaultTTL time.Duration
}

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

// NewQueryCache creates a new query cache with the specified default TTL
func NewQueryCache(defaultTTL time.Duration) *QueryCache {
	cache := &QueryCache{
		cache:      make(map[string]cacheEntry),
		defaultTTL: defaultTTL,
	}

	// Start a background goroutine to clean up expired entries
	go cache.cleanupLoop()

	return cache
}

// Get retrieves a value from the cache
func (c *QueryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	entry, found := c.cache[key]
	c.mutex.RUnlock()

	if !found {
		return nil, false
	}

	// Check if the entry has expired
	if time.Now().After(entry.expiration) {
		c.mutex.Lock()
		delete(c.cache, key)
		c.mutex.Unlock()
		return nil, false
	}

	return entry.value, true
}

// Set adds a value to the cache with the default TTL
func (c *QueryCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL adds a value to the cache with a specific TTL
func (c *QueryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *QueryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)
}

// Clear removes all values from the cache
func (c *QueryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]cacheEntry)
}

// cleanupLoop periodically removes expired entries from the cache
func (c *QueryCache) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for key, entry := range c.cache {
			if time.Now().After(entry.expiration) {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}

// Global query cache instance with a default TTL of 5 minutes
var QueryCacheInstance = NewQueryCache(5 * time.Second)

// QueryWithCache executes a database query with caching
// The queryFn parameter should be a function that executes the actual database query
// If the result is found in the cache, it is returned without executing the query
func QueryWithCache(cacheKey string, result interface{}, queryFn func() error) error {
	logger := appLog.WithFields(map[string]interface{}{
		"cache_key": cacheKey,
		"action":    "QueryWithCache",
	})

	// Try to get the result from the cache
	if cachedResult, found := QueryCacheInstance.Get(cacheKey); found {
		logger.Debug("Cache hit")

		// Copy the cached result to the result parameter
		switch typedResult := result.(type) {
		case *[]interface{}:
			*typedResult = cachedResult.([]interface{})
		case *interface{}:
			*typedResult = cachedResult
		default:
			// For other types, we need to use reflection or type assertions
			// This is a simplified implementation
			logger.Debug("Cache hit but result type not handled, executing query")
			goto executeQuery
		}

		return nil
	}

	logger.Debug("Cache miss, executing query")

executeQuery:
	// Execute the query
	err := queryFn()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Query execution failed")
		return err
	}

	// Store the result in the cache
	QueryCacheInstance.Set(cacheKey, result)
	logger.Debug("Query result cached")

	return nil
}
