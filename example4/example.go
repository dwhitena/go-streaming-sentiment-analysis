package main

import (
	"fmt"
	"sync"
	"time"
)

// Stats stores aggregated stats about
// tweets collected over time
type Stats struct {
	SentimentAverage float64
	Counts           map[string]int
	Mux              sync.Mutex
}

// IncrementCount increments the count of tweets.
func (s *Stats) IncrementCount(key string) {

	// Lock so only the current goroutine can access the map.
	s.Mux.Lock()

	// Increment the count.
	s.Counts[key]++

	// Unlock the data.
	s.Mux.Unlock()
}

// GetCount returns a count of tweets.
func (s *Stats) GetCount(key string) int {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	return s.Counts[key]
}

func main() {

	// Initialize our tweet stats.
	stats := &Stats{
		Counts: map[string]int{
			"positive": 0,
			"negative": 0,
			"neutral":  0,
		},
		Mux: sync.Mutex{},
	}

	for i := 0; i < 100; i++ {
		go stats.IncrementCount("positive")
	}

	time.Sleep(time.Second)
	fmt.Println(stats.GetCount("positive"))
}
