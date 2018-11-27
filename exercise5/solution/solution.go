package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/machinebox/sdk-go/textbox"
)

// Tweet is a single tweet.
type Tweet struct {
	Text  string
	Terms []string
}

// TweetReader includes the info we need to access Twitter.
type TweetReader struct {
	ConsumerKey, ConsumerSecret, AccessToken, AccessSecret string
}

// NewTweetReader creates a new TweetReader with the given credentials.
func NewTweetReader(consumerKey, consumerSecret, accessToken, accessSecret string) *TweetReader {
	return &TweetReader{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken:    accessToken,
		AccessSecret:   accessSecret,
	}
}

// Stats stores aggregated stats about
// tweets collected over time
type Stats struct {
	SentimentAverage float64
	Counts           map[string]int
	Mux              sync.Mutex
}

// IncrementCount increments the count of tweets.
func (s *Stats) IncrementCount(sentiment float64) {

	// Get the appropriate counter.
	var key string
	switch {
	case sentiment > 0.80:
		key = "positive"
	case sentiment < 0.50:
		key = "negative"
	default:
		key = "neutral"
	}

	// Update the counts.
	s.Mux.Lock()
	s.Counts[key]++
	s.Counts["total"]++
	s.Mux.Unlock()
}

// UpdateSentiment updates the tweet stream sentiment.
func (s *Stats) UpdateSentiment(newSentiment float64) {

	// Lock so only the current goroutine can access the sentiment.
	s.Mux.Lock()

	// Get the current count of tweets.
	total, ok := s.Counts["total"]
	if !ok {
		fmt.Println("Could not get key value \"total\"")
		return
	}

	// Update the value.
	s.SentimentAverage = (newSentiment + s.SentimentAverage*float64(total)) / (float64(total) + 1.0)

	// Unlock the data.
	s.Mux.Unlock()
}

func main() {

	// Create a new HTTP client.
	var connLock sync.Mutex
	var conn net.Conn
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				connLock.Lock()
				defer connLock.Unlock()
				if conn != nil {
					conn.Close()
					conn = nil
				}
				netc, err := net.DialTimeout(netw, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}
				conn = netc
				return netc, nil
			},
		},
	}

	// Create a new Tweet Reader.
	consumerKey := ""
	consumerSecret := ""
	accessToken := ""
	accessSecret := ""
	r := NewTweetReader(consumerKey, consumerSecret, accessToken, accessSecret)

	// Create oauth Credentials.
	creds := &oauth.Credentials{
		Token:  r.AccessToken,
		Secret: r.AccessSecret,
	}

	// Create an oauth Client.
	authClient := &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  r.ConsumerKey,
			Secret: r.ConsumerSecret,
		},
	}

	// Create the MachineBox client.
	machBoxIP := ""
	mbClient := textbox.New(machBoxIP)

	// Initialize the stats.
	myStats := Stats{
		SentimentAverage: 0.0,
		Counts: map[string]int{
			"positive": 0,
			"negative": 0,
			"neutral":  0,
			"total":    0,
		},
		Mux: sync.Mutex{},
	}

	// Setup the values we need for the context and filtering.
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	tweets := make(chan Tweet)
	terms := []string{"Trump", "Russia"}

	fmt.Println("Start 1st goroutine to collect tweets...")
	go func() {

		// Prepare the query.
		form := url.Values{"track": terms}
		formEnc := form.Encode()
		u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
		if err != nil {
			fmt.Println("Error parsing URL:", err)
		}

		// Prepare the request.
		req, err := http.NewRequest("POST", u.String(), strings.NewReader(formEnc))
		if err != nil {
			fmt.Println("creating filter request failed:", err)
		}
		req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", u, form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))

		// Execute the request.
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error getting response:", err)
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Unexpected HTTP status code:", resp.StatusCode)
		}

		// Decode the results.
		decoder := json.NewDecoder(resp.Body)
		for {
			var t Tweet
			if err := decoder.Decode(&t); err != nil {
				break
			}
			tweets <- t
		}
		resp.Body.Close()
	}()

	fmt.Println("Start a 2nd goroutine that prints the collected tweets...")
	go func() {
		for {
			select {

			// Stop the goroutine.
			case <-ctx.Done():
				return

			// Print the tweets.
			case t := <-tweets:

				// Analyze the tweet.
				analysis, err := mbClient.Check(strings.NewReader(t.Text))
				if err != nil {
					fmt.Println("MachineBox error:", err)
					continue
				}

				// Get the sentiment.
				sentimentTotal := 0.0
				for _, sentence := range analysis.Sentences {
					sentimentTotal += sentence.Sentiment
				}
				sentimentTotal = sentimentTotal / float64(len(analysis.Sentences))

				// Update the stats.
				myStats.UpdateSentiment(sentimentTotal)
				myStats.IncrementCount(sentimentTotal)
			}
		}
	}()

	// Check on our stats.
	for i := 0; i < 10; i++ {
		fmt.Println("")
		time.Sleep(time.Second)
		myStats.Mux.Lock()
		fmt.Printf("Sentiment: %0.2f\n", myStats.SentimentAverage)
		fmt.Printf("Total tweets analyzed: %d\n", myStats.Counts["total"])
		fmt.Printf("Total positive tweets: %d\n", myStats.Counts["positive"])
		fmt.Printf("Total negative tweets: %d\n", myStats.Counts["negative"])
		fmt.Printf("Total neutral tweets: %d\n", myStats.Counts["neutral"])
		myStats.Mux.Unlock()
	}

}
