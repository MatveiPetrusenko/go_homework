package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const limitValue = 100

func writer(ctx context.Context, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			close(ch)
			return
		default:
			ch <- rand.Intn(limitValue)
		}
	}
}

func reader(wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()

	for value := range ch {
		fmt.Printf("Recived %d\n", value)
	}
}

func main() {
	fmt.Println("Main Started")

	var flagN int
	flag.IntVar(&flagN, "flagN", 0, "An N flag")
	flag.Parse()

	var wg sync.WaitGroup
	ch := make(chan int, flagN)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(flagN)*time.Second)
	defer cancel()

	wg.Add(2)
	go writer(ctx, &wg, ch)
	go reader(&wg, ch)

	wg.Wait()
	fmt.Println("Main Ended")
}
