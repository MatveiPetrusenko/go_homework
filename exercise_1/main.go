package main

import (
	"flag"
	"fmt"
	"sync"
)

var nWorkers int

func main() {
	flag.IntVar(&nWorkers, "nWorkers", 0, "An N workers")
	flag.Parse()

	mainChannel := make(chan string)
	var wg sync.WaitGroup

	//start workers
	for i := 1; i <= nWorkers; i++ {
		wg.Add(1)
		go worker(i, mainChannel, &wg)
	}

	//write to main channel
	for i := 1; i <= nWorkers; i++ {
		data := fmt.Sprintf("Data %d", i)
		mainChannel <- data
	}
	close(mainChannel)

	wg.Wait()
}

func worker(id int, dataChannel <-chan string, wg *sync.WaitGroup) {
	for data := range dataChannel {
		fmt.Printf("Go_worker %d: Received %s\n", id, data)
	}

	defer wg.Done()
}
