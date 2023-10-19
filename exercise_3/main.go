package main

import (
	"fmt"
)

var array = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func main() {
	channelIn := make(chan int)
	channelOut := make(chan int)

	//go func writer
	go func(chan<- int) {
		for _, val := range array {
			channelIn <- val

			//fmt.Printf("Go_func:\"Writer\" send %d\n", val)
		}
		defer close(channelIn)
	}(channelIn)

	//go func double
	go func(<-chan int, chan<- int) {
		for val := range channelIn {
			channelOut <- val * val

			//fmt.Printf("Go_func:\"Double\" send %d\n", val*val)
		}

		defer close(channelOut)
	}(channelIn, channelOut)

	for val := range channelOut {
		fmt.Printf("Received %d from %T\n", val, channelOut)
	}
}
