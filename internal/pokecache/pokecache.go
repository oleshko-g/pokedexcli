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
	mu      sync.Mutex
	entries map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cacheEntry, ok := c.entries[key]
	return cacheEntry.val, ok
}

func (c *Cache) reapLoop(t *time.Ticker, interval time.Duration) {
	for tick := range t.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if tick.Sub(entry.createdAt) > interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	ticker := time.NewTicker(interval)
	cache := &Cache{
		mu:      sync.Mutex{},
		entries: make(map[string]cacheEntry),
	}
	go cache.reapLoop(ticker, interval)
	return cache
}
