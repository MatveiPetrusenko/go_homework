package main

import (
	"fmt"
	"math/rand"
)

const limitValue = 100

func generateSlice() []int {
	slice := make([]int, rand.Intn(limitValue))

	for k, _ := range slice {
		slice[k] = rand.Intn(limitValue)
	}

	return slice
}

func main() {
	slice := generateSlice()

	channelIn := make(chan int, limitValue)
	channelOut := make(chan int, limitValue)

	//go func writer
	go func(channelIn chan<- int) {
		defer close(channelIn)

		for _, val := range slice {
			channelIn <- val

			//fmt.Printf("Go_func:\"Writer\" send %d\n", val)
		}
	}(channelIn)

	//go func double
	go func(channelIn <-chan int, channelOut chan<- int) {
		defer close(channelOut)

		for val := range channelIn {
			channelOut <- val * val

			//fmt.Printf("Go_func:\"Double\" send %d\n", val*val)
		}
	}(channelIn, channelOut)

	for val := range channelOut {
		fmt.Printf("Received %d from %T\n", val, channelOut)
	}
}
