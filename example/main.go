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
	body, err := client.Classify("this movie is awesome")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", body)
}
