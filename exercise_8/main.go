package main

import (
	"fmt"
	lruCache "main.go/lru"
	"time"
)

func main() {
	cache := lruCache.NewLRUCache(lruCache.WithCapacity(1000), lruCache.WithTTL(1*time.Minute), lruCache.WithOtherOption())

	cache.Add("key1", "value1")

	res, ok := cache.Get("key1")
	fmt.Printf("%v %t\n", res, ok)

	res, notOk := cache.Get("key2")
	fmt.Printf("%v %t\n", res, notOk)

	time.Sleep(1 * time.Minute)

	res, notOk = cache.Get("key1")
	fmt.Printf("%v %t\n", res, notOk)
}
