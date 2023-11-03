package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os/signal"
	"sync"
	"syscall"
)

const limitValue = 100

func main() {
	var nWorkers int

	flag.IntVar(&nWorkers, "nWorkers", 1, "An N workers")
	flag.Parse()

	mainChannel := make(chan string, nWorkers)
	var wg sync.WaitGroup

	//create context which listening syscall + defer for out/stop from context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	//start workers
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go worker(i, mainChannel, &wg)
	}

	//write to main channel
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(mainChannel)
				return
			default:
				data := fmt.Sprintf("Data %d", rand.Intn(limitValue))
				mainChannel <- data
			}
		}
	}()

	wg.Wait()
}

func worker(id int, mainChannel <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range mainChannel {
		fmt.Printf("Go_worker %d: Received %s\n", id, data)
	}
}
