package cache

import "sync"

type Cache struct {
	mu   sync.Mutex
	data []interface{}
}

func New() *Cache {
	return &Cache{
		data: make([]interface{}, 0),
	}
}

func (c *Cache) Add(item interface{}) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := len(c.data)
	c.data = append(c.data, item)
	return id
}

func (c *Cache) GetAll() []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]interface{}, len(c.data))
	copy(out, c.data)
	return out
}
