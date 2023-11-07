package main

import (
	"fmt"
	"unicode/utf8"
)

const (
	left  = int32(40)
	right = int32(41)
)

func isBracketValid(s string) bool {
	if utf8.RuneCountInString(s) <= 1 {
		return false
	}

	var counter int

	for _, val := range s {
		if counter == -1 {
			return false
		}

		switch val {
		case left:
			counter++
		case right:
			counter--
		}
	}

	if counter == 0 {
		return true
	}

	return false
}

func main() {
	fmt.Printf("%t\n", isBracketValid("(())()")) //true
	fmt.Printf("%t\n", isBracketValid("())("))   //false
}
