package api

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	duration time.Duration
	cache    map[string]cacheEntry
	mu       *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	contents  []byte
}

func NewCache(dur time.Duration) Cache {
	mu := sync.RWMutex{}
	c := Cache{duration: dur, cache: map[string]cacheEntry{}, mu: &mu}
	go c.reapLoop(c.duration / 5)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{createdAt: time.Now(), contents: val}
}

func (c *Cache) Get(key string) (body []byte, presentInCache bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, pres := c.cache[key]
	return val.contents, pres
}

// Updates the entry's createdTime so that it staves off the reaper
func (c *Cache) Update(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.cache[key]; ok {
		entry.createdAt = time.Now()
		c.cache[key] = entry
	}
}

func (c *Cache) reapLoop(dur time.Duration) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()
	for t := range ticker.C {
		for k, v := range c.cache {
			// If the difference between now and the created time is greater than the duration variable
			if t.Sub(v.createdAt) >= c.duration {
				log.Printf("Culling entry : %s", k)
				c.mu.Lock()
				delete(c.cache, k)
				c.mu.Unlock()
			}
		}
	}
}
