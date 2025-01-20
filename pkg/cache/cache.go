package cache

import (
	"sync"
	"time"
)

type Item struct {
	Data      interface{}
	ExpiresAt time.Time
}

type Cache struct {
	store map[string]Item
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		store: make(map[string]Item),
		ttl:   ttl,
	}
}

// Get obtiene un elemento del cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.store[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Data, true
}

// Set agrega o actualiza un elemento en el cache.
func (c *Cache) Set(key string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = Item{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Delete elimina un elemento del cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}
