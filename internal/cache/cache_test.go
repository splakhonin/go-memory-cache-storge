package cache

import (
	"strconv"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache()

	value := "test value"
	key := "test-key"

	cache.Set(key, value)

	cachedValue, err := cache.Get(key)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if cachedValue != value {
		t.Errorf("Expected value '%s', but got '%v'", value, cachedValue)
	}
}

func TestCache_SetWithTTL(t *testing.T) {
	cache := NewCache()

	value := "test value"
	key := "test-key"
	ttl := 2 * time.Second

	cache.Set(key, value, ttl)

	cachedValue, err := cache.Get(key)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if cachedValue != value {
		t.Errorf("Expected value '%s', but got '%v'", value, cachedValue)
	}
}

func TestCache_GetWithTTL_Success(t *testing.T) {
	cache := NewCache()

	value := "test value"
	key := "test-key"
	ttl := 2 * time.Second

	cache.Set(key, value, ttl)

	// Sleep for slightly longer than TTL to ensure expiration
	time.Sleep(ttl + 1*time.Second)

	_, err := cache.Get(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCache_GetWithTTL_Failed(t *testing.T) {
	cache := NewCache()

	value := "test value"
	key := "test-key"
	ttl := 2 * time.Second

	cache.Set(key, value, ttl)

	// Sleep for slightly longer than TTL to ensure expiration
	time.Sleep(ttl + 4*time.Second)

	_, err := cache.Get(key)
	if err != nil {
		// No assertion needed, test passes if an error is present
		return
	}

	t.Error("Expected an error due to TTL expiration, but got NO error")
}

func TestCache_CacheClean(t *testing.T) {
	cache := NewCache()

	value := "test value"
	key := "test-key"

	cache.Set(key, value)

	err := cache.CacheClean(key)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	_, err = cache.Get(key)
	if err == nil {
		t.Errorf("Expected an error after cleaning cache, but got no error")
	}
}

func TestCache_StopCleaner(t *testing.T) {
	cache := NewCache()

	// Add some values to cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Wait for a moment to ensure the cleaner has a chance to run
	time.Sleep(1 * time.Second)

	// Stop the cleaner
	cache.StopCleaner()

	// Wait for a moment to ensure the cleaner is fully stopped
	time.Sleep(1 * time.Second)

	// Try to clean cache manually (which shouldn't work after stopping the cleaner)
	cache.clean()

	// Ensure that the cache is not cleaned and still contains the old values
	_, err := cache.Get("key1")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	_, err = cache.Get("key2")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache()

	// Concurrently set and get values
	const numRoutines = 10
	for i := 0; i < numRoutines; i++ {
		go func(idx int) {
			key := "key" + strconv.Itoa(idx)
			value := "value" + strconv.Itoa(idx)
			cache.Set(key, value)
			cachedValue, err := cache.Get(key)
			if err != nil {
				t.Errorf("Routine %d: Expected no error, but got: %v", idx, err)
			}
			if cachedValue != value {
				t.Errorf("Routine %d: Expected value '%s', but got '%v'", idx, value, cachedValue)
			}
		}(i)
	}

	// Wait for routines to complete
	time.Sleep(2 * time.Second)
}
