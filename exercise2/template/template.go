package main

import (
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

	// Create a new HTTP client. Because we are going to do some streaming analysis
	// and may be utilizing this client from mutliple goroutines, I have borrowed
	// the HTTP client configuration from a similar MachineBox project that will keep
	// us safe in these scenarios. It includes both a Mutex for accessing the client
	// and some timeout functionality.
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

	// TODO: Create a new Tweet Reader (Fill in your Twitter keys here)
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

	// TODO: Define the terms for our search. Create a slice of strings
	// value with the terms you want to search.

	// Form the URL for our search.
	form := url.Values{"track": terms}
	formEnc := form.Encode()

	// Create a new HTTP request.
	u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	if err != nil {
		fmt.Println("Could not parse url:", err)
	}

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(formEnc))
	if err != nil {
		fmt.Println("creating filter request failed:", err)
	}

	// Set some header info.
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", u, form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))

	// Do the request.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting response:", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("failed with status code:", resp.StatusCode)
	}

	// TODO: Create a new decoder for the response body. Hint, use the
	// "encoding/json" package and call the decoder "decoder".

	// Start reading in tweets and parsing them.
	for i := 0; i < 10; i++ {
		var t Tweet
		if err := decoder.Decode(&t); err != nil {
			break
		}
		fmt.Printf("TWEET %d TEXT: %s\n", i+1, t.Text)
		fmt.Println("----------------------------------------\n")
	}

	// Close the response body.
	resp.Body.Close()
}
