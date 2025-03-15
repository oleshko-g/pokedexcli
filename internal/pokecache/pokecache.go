package pokecache

import (
	"sync"
)

type cacheEntry struct {
	val []byte
}

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		val: val,
	}
}

func (c *Cache) Get(key string) (val []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cacheEntry, ok := c.entries[key]
	return cacheEntry.val, ok
}

func NewCache() Cache {
	return Cache{
		mu:      sync.Mutex{},
		entries: make(map[string]cacheEntry),
	}
}
