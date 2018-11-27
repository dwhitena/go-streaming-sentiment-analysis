package main

import (
	"fmt"
	"time"
)

func main() {

	// Create our buffered channel.
	workers := 20
	results := make(chan string, workers)

	// Have our workers complete some expensive tasks.
	for w := 0; w < workers; w++ {
		output := fmt.Sprintf("Worker %d output", w)
		go func(output string) {
			time.Sleep(1 * time.Second)
			results <- output
		}(output)
	}

	// Collect the results.
	for i := 0; i < workers; i++ {
		p := <-results
		fmt.Println(p)
	}
}
