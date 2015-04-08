package main

import (
	"fmt"
	"log"
	"os"

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
	fmt.Printf("Sentiment for: %s\n", testText)
	for _, p := range sent.Predictions {
		fmt.Printf("%s %g\n", p.ClassName, p.Prob)
	}
}
