package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Lock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Lock()
	cacheEntry, ok := c.entries[key]
	return cacheEntry.val, ok
}

//func (c *Cache) reapCache(interval time.Duration) {
//	for _, v := range c.entries {
//		if time.Now().k - v.createdAt > interval {
//
//		}
//	}
//}

func NewCache(interval time.Duration) (Cache, error) {
	return Cache{
		entries: make(map[string]cacheEntry),
		mu:      sync.Mutex{},
	}, nil
}
