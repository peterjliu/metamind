package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/peterjliu/metamind"
)

type SentimentResp struct {
	Text       string
	Prediction *metamind.SentiResp
	Success    bool
}

func main() {
	metamind_key := os.Getenv("METAMIND_API")
	client := metamind.NewSentimentClient(metamind_key)
	testText := "this movie is awesome"
	sent, err := client.Classify(testText)
	if err != nil {
		log.Fatal(err)
	}
	printSentiments(testText, sent)

	// Tweets
	consumerkey := os.Getenv("TWITTER_C")
	consumersecret := os.Getenv("TWITTER_CS")
	accesstoken := os.Getenv("TWITTER_AT")
	accesstokensecret := os.Getenv("TWITTER_ST")

	// fmt.Println("%s %s %s %s\n", consumerkey, consumersecret, accesstoken, accesstokensecret)
	anaconda.SetConsumerKey(consumerkey)
	anaconda.SetConsumerSecret(consumersecret)
	api := anaconda.NewTwitterApi(accesstoken, accesstokensecret)

	sentiments := make(chan *SentimentResp)
	searchResult, _ := api.GetSearch("golang", nil)
	for _, tweet := range searchResult.Statuses {
		go GetSentiment(client, tweet.Text, sentiments)
	}
	num_failed := 0
	for i := 0; i < len(searchResult.Statuses); i++ {
		r := <-sentiments
		if r.Success {
			printSentiments(r.Text, r.Prediction)
		} else {
			num_failed += 1
		}
	}
}

func printSentiments(testText string, senti *metamind.SentiResp) {
	fmt.Printf("Sentiment for: %s\n", testText)
	for _, p := range senti.Predictions {
		fmt.Printf("%s %g\n", p.ClassName, p.Prob)
	}
	fmt.Printf("\n")
}

func GetSentiment(client *metamind.SentiClient, tweet string, sentiments chan *SentimentResp) {
	senti, err := client.Classify(tweet)
	sentiments <- &SentimentResp{Text: tweet, Prediction: senti, Success: err == nil}
}
