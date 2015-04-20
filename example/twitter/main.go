package main

import (
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/peterjliu/metamind"
)

func main() {
	metamind_key := os.Getenv("METAMIND_API")
	client := metamind.NewSentimentClient(metamind_key)
	testText := "this movie is awesome"
	sent, err := client.Classify(testText)
	if err != nil {
		log.Fatal(err)
	}
	metamind.PrintSentiments(testText, sent)

	// Tweets
	consumerkey := os.Getenv("TWITTER_C")
	consumersecret := os.Getenv("TWITTER_CS")
	accesstoken := os.Getenv("TWITTER_AT")
	accesstokensecret := os.Getenv("TWITTER_ST")

	// fmt.Println("%s %s %s %s\n", consumerkey, consumersecret, accesstoken, accesstokensecret)
	anaconda.SetConsumerKey(consumerkey)
	anaconda.SetConsumerSecret(consumersecret)
	api := anaconda.NewTwitterApi(accesstoken, accesstokensecret)

	sentiments := make(chan *metamind.SentimentResp)
	searchResult, _ := api.GetSearch("golang", nil)
	for _, tweet := range searchResult.Statuses {
		go metamind.GetSentiment(client, tweet.Text, sentiments)
	}
	num_failed := 0
	for i := 0; i < len(searchResult.Statuses); i++ {
		r := <-sentiments
		if r.Success {
			metamind.PrintSentiments(r.Text, r.Prediction)
		} else {
			num_failed += 1
		}
	}
}
