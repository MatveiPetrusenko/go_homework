## N.8
> Реализовать универсальный **LRU-cache**\
Регулируется **TTL** и размер кеша. Реализовать в виде **go module**.\
Добавить тесты, в т.ч. тестирование производительности и потребления памяти.

```
func main() {
    cache := NewLRUCache(WithCapacity(1000), WithTTL(1 * time.Minute), WithOtherOption())
    cache.Add("key1", value)

    res, ok := cache.Get("key1")
    res, notOk := cache.Get("key2")

    time.Sleep(1 * time.Minute)

    res, notOk := cache.Get("key1")
}
```


