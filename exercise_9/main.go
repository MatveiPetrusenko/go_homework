package main

import (
	"fmt"
)

func isBracketValid(s string) bool {
	if s[0] == ')' || s[0] == ']' || s[0] == '}' {
		return false
	}

	brackets := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	stack := make([]rune, 0, len(s))
	for _, val := range s {
		switch val {
		case '(', '[', '{':
			stack = append(stack, val)
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != brackets[val] {
				return false
			}

			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func main() {
	fmt.Printf("%t\n", isBracketValid("(())()")) //true
	fmt.Printf("%t\n", isBracketValid("())("))   //false
	fmt.Printf("%t\n", isBracketValid(")())("))  //false
}
