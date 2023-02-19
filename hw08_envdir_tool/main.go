package main

import (
	"fmt"
)

func main() {
	rez, err := ReadDir("./testdata/env")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("general_out: ", rez)
}
