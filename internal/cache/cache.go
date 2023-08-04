package cache

import "errors"

type Cache struct {
	storage map[string]interface{}
}

// NewCache initialize a new cache instance
func NewCache() *Cache {
	return &Cache{
		storage: make(map[string]interface{}),
	}
}

// Set a new value to cache
func (c *Cache) Set(key string, value interface{}) {
	c.storage[key] = value
}

// Get a value from cache
func (c *Cache) Get(key string) (interface{}, error) {
	val, ok := c.storage[key]
	if !ok {
		return nil, errors.New("can not get a key, it does not exist")
	}

	return val, nil
}

// Delete a value from cache
func (c *Cache) Delete(key string) error {
	_, ok := c.storage[key]
	if !ok {
		return errors.New("can not delete by key, it does not exist")
	}
	delete(c.storage, key)
	return nil
}
