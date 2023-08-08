package cache

import (
	"errors"
	"sync"
	"time"
)

const (
	cleanInterval = 5 * time.Second
	defaultTTL    = 5 * time.Second // Default TTL
)

type Cache struct {
	storage  map[string]cachedData
	mu       *sync.RWMutex
	cleaner  *time.Ticker
	stopChan chan struct{}
}

type cachedData struct {
	value  interface{}
	expiry time.Time
}

// NewCache initializes a new cache instance with a background cleaner
func NewCache() *Cache {
	c := &Cache{
		storage:  make(map[string]cachedData),
		mu:       &sync.RWMutex{},
		cleaner:  time.NewTicker(cleanInterval),
		stopChan: make(chan struct{}),
	}

	go c.startCleaner()

	return c
}

// Clear a value from cache
func (c *Cache) startCleaner() {
	for {
		select {
		case <-c.cleaner.C:
			c.clean()
		case <-c.stopChan:
			c.cleaner.Stop()
			return
		}
	}
}

// StopCleaner stop the background cleaner
func (c *Cache) StopCleaner() {
	close(c.stopChan)
}

func (c *Cache) clean() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, cacheData := range c.storage {
		if cacheData.expiry.Before(now) {
			delete(c.storage, key)
		}
	}
}

// Set a new value to cache with an optional TTL
func (c *Cache) Set(key string, value interface{}, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiry time.Time

	// If a TTL is provided, use it; otherwise, use the default
	if len(ttl) > 0 {
		expiry = time.Now().Add(ttl[0])
	} else {
		expiry = time.Now().Add(defaultTTL) // Use default TTL from the constant
	}

	c.storage[key] = cachedData{
		value:  value,
		expiry: expiry,
	}
}

// Get a value from cache
func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cacheData, ok := c.storage[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}

	return cacheData.value, nil
}

// CacheClean removes a value from cache by key
func (c *Cache) CacheClean(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.storage[key]
	if !ok {
		return errors.New("key does not exist")
	}

	delete(c.storage, key)
	return nil
}
