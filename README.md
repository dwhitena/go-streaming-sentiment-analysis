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

This repo also include [a number of resources (at the end of this README)](#resources) if you want to dig into more details.

## 1. Install Go

You will be running all of the tutorial steps on your local machine. Thus, you need to make sure that you have Go installed. I recommend following [this guide for setting up your local Go environment](https://www.ardanlabs.com/blog/2016/05/installing-go-and-your-workspace.html). 

*Note - I recommend that you do NOT use a package manager like `brew` to install Go. These tend to make it harder to manage your Go environment (especially with respect to updates)*

## 2. Run MachineBox's Text Box

[MachineBox](https://machinebox.io/) is state-of-the-art machine learning technology inside a Docker container which you can run, deploy and scale. We will utilize "textbox," which is a pre-configured Docker image for sentiment analysis. Everyone can run their own textbox on their own machine as follows:

1. Open a new terminal window
2. Go to [machinebox.io](https://machinebox.io/) and click "sign up for your free key"
3. Enter your email address
4. Check for the sign up email in your inbox
5. Follow the link in the email. This will take you to your account page on the MachineBox site.
6. Scroll down to the bottom of that account page and click on "textbox"
7. Follow the instructions to run `textbox` via Docker

## Resources 

blah
