package lru

import (
	"container/list"
	"sync"
	"time"
)

type LRUCache struct {
	mu       sync.Mutex
	capacity int
	ll       *list.List
	ttl      time.Duration
	cache    map[interface{}]*list.Element
}

type cacheEntry struct {
	key        interface{}
	value      interface{}
	expiration time.Time
}

type Option func(*LRUCache)

func WithCapacity(capacity int) Option {
	return func(c *LRUCache) {
		c.capacity = capacity
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(c *LRUCache) {
		c.ttl = ttl
	}
}

func WithOtherOption() Option {
	return func(c *LRUCache) {
	}
}

func NewLRUCache(options ...Option) *LRUCache {
	newCache := &LRUCache{
		ll:    list.New(),
		cache: make(map[interface{}]*list.Element),
	}

	for _, option := range options {
		option(newCache)
	}

	return newCache
}

func (c *LRUCache) Add(key interface{}, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)

		entry := element.Value.(*cacheEntry)
		entry.value = value
		entry.expiration = time.Now().Add(c.ttl)
	} else {
		if len(c.cache) >= c.capacity {
			c.expired()
		}

		entry := &cacheEntry{
			key:        key,
			value:      value,
			expiration: time.Now().Add(c.ttl),
		}
		element := c.ll.PushFront(entry)
		c.cache[key] = element
	}
}

func (c *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.cache[key]; ok {
		entry := element.Value.(*cacheEntry)
		if time.Now().Before(entry.expiration) {
			c.ll.MoveToFront(element)

			return entry.value, true
		}

		c.remove(element)
	}

	return nil, false
}

func (c *LRUCache) remove(element *list.Element) {
	entry := element.Value.(*cacheEntry)
	delete(c.cache, entry.key)
	c.ll.Remove(element)
}

func (c *LRUCache) expired() {
	element := c.ll.Back()
	if element != nil {
		c.remove(element)
	}
}
