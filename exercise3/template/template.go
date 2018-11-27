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

	// Create a "tweets" channel.
	tweets := make(chan Tweet)

	// Define the terms for our search.
	terms := []string{"Trump", "Russia"}

	// Create a context value that will allow us to stop our goroutine.
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

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

			// TODO: Push the tweet value t to the tweets channel.
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

				// TODO: Create a case statement that reads tweets off of the tweets
				// channel and prints them.
			}
		}
	}()

	time.Sleep(3 * time.Second)
}
