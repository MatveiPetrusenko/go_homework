package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var nWorkers int
var limitValue = 100

func main() {
	flag.IntVar(&nWorkers, "nWorkers", 0, "An N workers")
	flag.Parse()

	mainChannel := make(chan string, nWorkers)
	var wg sync.WaitGroup

	//create signal channel + catcher for syscall
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

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

	<-signals
	wg.Wait()
}

func worker(id int, mainChannel <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range mainChannel {
		fmt.Printf("Go_worker %d: Received %s\n", id, data)
	}
}
