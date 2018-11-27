package main

import (
	"fmt"
	"time"
)

// printNum is a function that prints out numbers.
func printNum(from string) {
	for i := 1; i < 5; i++ {
		fmt.Printf("%s output: %d\n", from, i)
		time.Sleep(1 * time.Second)
	}
}

func main() {

	// Start a goroutine that will print some numbers.
	go printNum("goroutine 1")

	// Start another goroutine that will print some numbers.
	go func() {
		fmt.Println("goroutine 2 output!")
	}()

	time.Sleep(4 * time.Second)
}
