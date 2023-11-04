package lru

import (
	lru "github.com/hashicorp/golang-lru"
	"time"
)

type LRUCache struct {
	cache    *lru.Cache
	ttl      time.Duration
	ttlCache map[interface{}]int64
}

func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	lruCache, _ := lru.New(capacity)

	return &LRUCache{
		cache:    lruCache,
		ttl:      ttl,
		ttlCache: make(map[interface{}]int64),
	}
}

func (c *LRUCache) Add(key interface{}, value interface{}) {
	c.cache.Add(key, value)
	c.ttlCache[key] = time.Now().Unix()
}

func (c *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	value, ok = c.cache.Get(key)
	if ok {
		if c.isExpired(key) {
			c.Remove(key)
			return nil, false
		}
	}
	return value, ok
}

func (c *LRUCache) Remove(key interface{}) {
	c.cache.Remove(key)
	delete(c.ttlCache, key)
}

func (c *LRUCache) isExpired(key interface{}) bool {
	if c.ttl == 0 {
		return false
	}
	expirationTime := c.ttlCache[key] + int64(c.ttl.Seconds())
	currentTime := time.Now().Unix()
	return currentTime > expirationTime
}
