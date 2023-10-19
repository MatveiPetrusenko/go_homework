package main

import (
	"flag"
	"fmt"
	"time"
)

var flagN int

func reader(ch <-chan int) {
	for value := range ch {
		fmt.Printf("Recived %d from %v\n", value, flagN)
	}
}

func writer(ch chan<- int) {
	for i := 1; i <= flagN; i++ {
		ch <- i
		time.Sleep(time.Second)
	}

	close(ch)
}

func main() {
	fmt.Println("Main Started")

	flag.IntVar(&flagN, "flagN", 0, "An N flag")
	flag.Parse()

	ch := make(chan int)

	go writer(ch)
	go reader(ch)

	time.Sleep(time.Second * time.Duration(flagN))
	fmt.Println("Main Ended")
}
