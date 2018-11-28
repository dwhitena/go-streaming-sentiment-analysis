package main

import (
	"fmt"

	"github.com/machinebox/sdk-go/textbox"
)

func main() {

	// Define the MachineBox IP for our text box.
	machBoxIP := "http://localhost:8080"

	// Create a new MachineBox client.
	client := textbox.New(machBoxIP)

	// TODO: Define a positive or negative sentiment statements as a string.
	negativeStatement := "It is sad, depressing, and unfortunate that this workshop will terminate at the end of the day."

	// TODO: Try to get the sentiment of the positive and/or negative statement using
	// the MachineBox Check method.

	// Print out the keywords returned by MachineBox.
	fmt.Println(analysis.Keywords)

	// Compute the total sentiment.
	sentimentTotal := 0.0
	for _, sentence := range analysis.Sentences {
		sentimentTotal += sentence.Sentiment
	}

	// Print out the average sentiment. Higher sentitment is more positive,
	// and lower is more negative.
	fmt.Println("Sentiment:", sentimentTotal/float64(len(analysis.Sentences)))
}
