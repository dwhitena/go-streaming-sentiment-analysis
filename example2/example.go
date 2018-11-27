package main

import (
	"fmt"
	"time"
)

func main() {

	// Make a channel that will signal with string data.
	ch := make(chan string)

	// Start a goroutine that will print received string messages.
	go func() {
		p := <-ch
		fmt.Println("received signal: ", p)
	}()

	time.Sleep(time.Second)

	// Send a message on the channel.
	ch <- "This is my message"
	fmt.Println("sent signal")

	time.Sleep(time.Second)
}
