package main

import (
	"fmt"
	"strings"

	"github.com/machinebox/sdk-go/textbox"
)

func main() {

	// Define the MachineBox IP for our text box.
	machBoxIP := ""

	// Create a new MachineBox client.
	client := textbox.New(machBoxIP)

	// Define positive and negative sentiment statements.
	positiveStatement := "I am so excited to be teaching to super awesome, fun workshop!"
	negativeStatement := "It is sad, depressing, and unfortunate that this workshop will terminate at the end of the day."

	// Try to get the sentiment of the positive and/or negative statement using
	// the MachineBox Check method.
	analysis, err := client.Check(strings.NewReader(negativeStatement))
	if err != nil {
		fmt.Println(err)
	}

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
