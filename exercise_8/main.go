package main

import (
	"fmt"
	lruCache "main.go/lru"
	"time"
)

func main() {
	cache := lruCache.NewLRUCache(1000, 1*time.Minute)

	cache.Add("key1", "value1")
	res, ok := cache.Get("key1")
	fmt.Println(res, ok)

	res, notOk := cache.Get("key2")
	fmt.Println(res, notOk)

	time.Sleep(1 * time.Minute)

	res, notOk = cache.Get("key1")
	fmt.Println(res, notOk)
}
