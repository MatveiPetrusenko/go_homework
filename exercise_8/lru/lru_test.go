package lru

import (
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(2, 1*time.Minute)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")

	value, ok := cache.Get("key1")
	if !ok || value != "value1" {
		t.Errorf("Expected key1 to be present in the cache")
	}

	value, ok = cache.Get("key2")
	if !ok || value != "value2" {
		t.Errorf("Expected key2 to be present in the cache")
	}

	cache.Add("key3", "value3")

	value, ok = cache.Get("key1")
	if ok {
		t.Errorf("Expected key1 to be evicted from the cache")
	}
}

func BenchmarkAdd(b *testing.B) {
	cache := NewLRUCache(1000, 0)

	for i := 0; i < b.N; i++ {
		cache.Add(i, i)
	}
}

func BenchmarkGet(b *testing.B) {
	cache := NewLRUCache(1000, 0)
	for i := 0; i < 1000; i++ {
		cache.Add(i, i)
	}

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(i % 1000)
	}
}
