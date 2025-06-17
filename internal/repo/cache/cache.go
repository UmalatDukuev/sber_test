// Package cache provides a simple in-memory cache.
package cache

import "sync"

// Cache stores data in memory with thread-safe operations.
type Cache struct {
	data []interface{}
	mu   sync.Mutex
}

// New creates a new Cache instance.
func New() *Cache {
	return &Cache{
		data: make([]interface{}, 0),
	}
}

// Add adds an item to the cache and returns its ID.
func (c *Cache) Add(item interface{}) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := len(c.data)
	c.data = append(c.data, item)
	return id
}

// GetAll retrieves all items from the cache.
func (c *Cache) GetAll() []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]interface{}, len(c.data))
	copy(out, c.data)
	return out
}
