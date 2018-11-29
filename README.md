# Streaming Sentiment Analysis with MachineBox and Go

This repo will walk you through a tutorial for streaming sentiment analysis of tweets using [Go](https://golang.org/) and [MachineBox](https://machinebox.io/). By the end of the tutorial you will be able to retrieve tweets, calculate a ["sentiment"](https://en.wikipedia.org/wiki/Sentiment_analysis) of the tweets, and maintain aggregate statistics for certain topics on Twitter.

Prerequisites:

- A laptop
- WiFi
- [Docker](https://docs.docker.com/install/)
- A Twitter account

You can complete the tutorial by walking through the following steps:

1. [Install Go]()
2. [Run MachineBox's Text Box]()
3. [Create a Twitter "app"]()
4. [Learn how to work with and write Go]()
5. [Learn about a few relevant Go primitives]()
6. [Process tweets with Go]()
7. [Analyze sentiment with Go and MachineBox]()
8. [Analyze the sentiment of streaming tweets]()

## 1. Install Go

You will be running all of the tutorial steps on your local machine. Thus, you need to make sure that you have Go installed. I recommend following [this guide for setting up your local Go environment](https://www.ardanlabs.com/blog/2016/05/installing-go-and-your-workspace.html). 

*Note - I recommend that you do NOT use a package manager like `brew` to install Go. These tend to make it harder to manage your Go environment (especially with respect to updates)*

After installing Go, create the directory `$HOME/go/src` (if it doesn't already exist). Navigate to this `src` directory, then clone the tutorial materials:

```
$ git clone https://github.com/dwhitena/go-streaming-sentiment-analysis.git
```

## 2. Run MachineBox's Text Box

[MachineBox](https://machinebox.io/) is state-of-the-art machine learning technology inside a Docker container which you can run, deploy and scale. We will utilize "textbox," which is a pre-configured Docker image for sentiment analysis. Everyone can run their own textbox on their own machine as follows:

1. Open a new terminal window
2. Go to [machinebox.io](https://machinebox.io/) and click "sign up for your free key"
3. Enter your email address
4. Check for the sign up email in your inbox
5. Follow the link in the email. This will take you to your account page on the MachineBox site.
6. Scroll down to the bottom of that account page and click on "textbox"
7. Follow the instructions to run `textbox` via Docker

To confirm MachineBox is up and running you should be able to visit [http://localhost:8080/](http://localhost:8080/) and see the web interface to the textbox.

## 3. Create a Twitter "app"

To retrieve Tweets from Twitter's streaming API, you will need to obtain a set of credentials from your Twitter account. To do this:

1. Create an account at [developer.twitter.com](https://developer.twitter.com/)
2. Create a new twitter "app"
3. Under that app, under Keys and Access Tokens retrieve your connection key and secret.
4. Generate an access token and token secret on that same page.
5. Retrieve the access token and token secret.

## 4. Learn how to work with and write Go

The Go language provides an excellent platform for real-time analysis. To get your feet wet writing and running Go programs (and actually for all sections of this tutorial), we have created a series of exercies. Each exercise includes a template code file and a solution code file. In order to complete the exercise, first open up the template code file in your editor of choice. Then find the sections with `TODO` comments and complete the code required for those sections. Try to run your solution (based on the template file) by executing the following from the directory containing the template file:

```
$ go build
$ ./template
``` 

(this will build the binary for your program and then execute it). To find some resources and example of Go syntax and language primitives, try referencing [Go by example](https://gobyexample.com/). This should give you some ideas as you try to complete the exercises. 

Once you try out the exercises based on the template files, feel free to look at the corresponding solution file. It's totally ok to look at this if you are stuck! Also, note that these solutions just represent one way of solving the problem. There are enumerable other (completely valid) ways.

**Exercise 1** - This first exercise will introduce you to Go structs and how to use them to store data. From the root of your cloned version of this repo (in your `$HOME/go/src` directory), navigate to the `exercise1` directory. You will find the template and solution files for this exercise there.

In this exercise, we will create a `Tweet` struct that will hold the contents of an individual Tweet. In addition to the text of the Tweet, we will include a list of keywords that were matched in the tweet (as we will be searching through tweets by keyword).  We are also going to utilize a `Stats` struct that will allow us to aggregate some streaming statistics about tweets that we are analyzing.  The `SentimentAverage` field will hold a float value that will represent the current average sentiment (positive/negative or happy/sad) of analyzed tweets. This number will fluctuate between 1 and 0, but we will get into those details later. The `Counts` field of the `Stats` struct will include a map that we will update with counts of positive, negative, and neutral tweets.

## Learn about a few relevant Go primitives

We will utilize a number of unique Go language features to build our streaming Tweet analyze. In particular, we will utilize channels, goroutines, and mutexes. Before moving forward, let's learn a little more about these things.

### goroutines

A goroutine is a lightweight thread of execution. In the words of Bill Kennedy:

> Goroutines are functions that are created and scheduled to be run independently by the Go scheduler. The Go scheduler is responsible for the management and execution of goroutines.

These goroutines aren't the same as OS "threads." Threads are expensive OS operations. You run out of system resources very quickly when launching threads. Goroutines are more like a "coroutine" or "coprogram." These coroutines multiplex independently executing functions onto a set of threads.

Review this [goroutine example](example1)

### channels

Go channels facilitate communication or signaling between goroutines. We will cover some of the basic patterns of channel usage here, but there are a whole variety of types and patterns for Go channels. I highly recommend that you read [this blog post](https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html) for a detailed introduction.

Review this [unbuffered channel example](example2)

When signaling with some kind of data (e.g., strings as above), we can use a buffered or unbuffered channel. Our choice will have some impacts on whether we can garauntee delivery of the signal.

- If you need to have a goroutine do something only when it receives a signal (i.e., it needs to wait for a signal), you should use an unbuffered channel because you need to ensure that the signal is received.
- When you need to throw multiple workers at a problem and don't need to ensure receipt of an individual signal, you might use a buffered channel with a well defined number of workers. You might also use a buffered channel when you need to drop signals after being saturated with a well defined number of signals.

The example included above in the "Channels" section is an example of of the first case, when a goroutine is waiting on a signal. Once the goroutine receives the signal, it does the corresponding work (printing the message in this case).

Review this [buffered channel example](example3)

### mutexes

One of the main mechanisms for managing state in Go is communication over channels, as we learned in the previous notebook. However, as stated in [this useful article](https://github.com/golang/go/wiki/MutexOrChannel) in the official Go wiki:

> A common Go newbie mistake is to over-use channels and goroutines just because it's possible, and/or because it's fun. Don't be afraid to use a sync.Mutex if that fits your problem best.

So what does that mean, and what the heck is a "Mutex"? Well, a mutex is a mutual exclusion lock that can be utilized as a rule, such that we can safely access data across multiple goroutines. More specifically, a Mutex allows us to Lock certain data (e.g., a struct) that may also be accessed by other goroutines, such that we ensure exclusive access to the data until we Unlock that data.

Review this [mutex example](example4)

## Process tweets with Go

**Exercise 2** - This first exercise will introduce you to how we can process Tweets with Go. From the root of your cloned version of this repo (in your `$HOME/go/src` directory), navigate to the `exercise2` directory. You will find the template and solution files for this exercise there.

## Analyze sentiment with Go and MachineBox

**Exercise 3** - This third exercise will introduce you to how we can perform sentiment analysis with MachineBox from Go. Note, that you will need to have the textbox up and running in Docker as detailed above. From the root of your cloned version of this repo (in your `$HOME/go/src` directory), navigate to the `exercise3` directory. You will find the template and solution files for this exercise there.

## Analyze the sentiment of streaming tweets

**Exercise 4** - This third exercise will introduce you to how we can connect all of the building blocks and perform streaming sentiment analysis of tweets. Note, that you will need to have the textbox up and running in Docker as detailed above. From the root of your cloned version of this repo (in your `$HOME/go/src` directory), navigate to the `exercise4` directory. You will find the template and solution files for this exercise there.
