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

	// Initialize a value of Tweet.
	var tweet Tweet

	// Fill in some details.
	tweet.Text = "This is my really happy tweet!"
	tweet.Terms = []string{"tweet"}

	// Initialize and update our tweet stats.
	stats := Stats{
		SentimentAverage: 0.9,
		Counts: map[string]int{
			"positive": 1,
			"negative": 0,
			"neutral":  0,
		},
	}

	fmt.Printf("We have %d positive tweet(s) and %d negative tweet(s).\n",
		stats.Counts["positive"], stats.Counts["negative"])
}
