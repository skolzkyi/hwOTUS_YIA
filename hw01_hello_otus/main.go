package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	string := "Hello, OTUS!"
	reversedString := reverseString(string)
	fmt.Println(reversedString)
}

func reverseString(input string) string {
	output := stringutil.Reverse(input)
	return output
}
