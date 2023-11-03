package main

import (
	"fmt"
	"sync"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	single := make(chan interface{})

	for _, ch := range channels {
		wg.Add(1)

		go func(c <-chan interface{}) {
			defer wg.Done()

			/*			for _ = range c {
						}
						single <- struct{}{}*/

			for {
				val, ok := <-c
				if ok {
					continue
				}
				single <- val
				break
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(single)
	}()

	return single
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(3*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %s\n", time.Since(start))
}
