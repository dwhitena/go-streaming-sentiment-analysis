package main

import "fmt"

// Tweet is a single tweet.
type Tweet struct {
	Text  string
	Terms []string
}

// Stats stores aggregated stats about
// tweets collected over time
type Stats struct {
	SentimentAverage float64
	Counts           map[string]int
}

func main() {

	// TODO: Initialize a value of Tweet (hint - use the keyword "var").

	// TODO: Fill in the "Text" field of the tweet with some text
	// and fill in the Terms in the Tweet as a slice of strings.

	// TODO: Initialize and update our tweet stats. Update the numbers
	// below based on whether you think your tweet is positive
	// or negative.
	stats := Stats{
		SentimentAverage: 0.0,
		Counts: map[string]int{
			"positive": 0,
			"negative": 0,
			"neutral":  0,
		},
	}

	fmt.Printf("We have %d positive tweet(s) and %d negative tweet(s).\n",
		stats.Counts["positive"], stats.Counts["negative"])
}
