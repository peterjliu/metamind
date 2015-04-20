package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"

	"github.com/peterjliu/metamind"
	"github.com/peterjliu/textproc"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var text_file = flag.String("infile", "", "some file with text")
var limit = flag.Int("limit", 100, "Max limit on number of sentences")
var max_reqs = flag.Int("max_reqs", 10, "maximum outstanding requests")

func main() {
	flag.Parse()
	log.Printf("Reading from %s", *text_file)
	data, err := ioutil.ReadFile(*text_file)
	check(err)
	sentences := textproc.GetSentences(string(data))
	log.Printf("Got %d sentences", len(sentences))

	metamind_key := os.Getenv("METAMIND_API")
	client := metamind.NewSentimentClient(metamind_key)
	sentiments := make(chan *metamind.SentimentResp)
	// set a maximum outstanding requests
	sem := make(chan int, *max_reqs)
	t0 := time.Now()
	for i, s := range sentences {
		if i >= *limit {
			break
		}
		// Launch workers
		go func(s2 string) {
			sem <- 1
			err := metamind.GetSentiment(client, s2, sentiments)
			check(err)
		}(s)

	}
	num_failed := 0
	// Start channel consumer
	total_req_time_ms := 0.0
	max_req_time_ms := -1.0
	min_req_time_ms := math.MaxFloat32
	for i := 0; i < *limit; i++ {
		<-sem
		r := <-sentiments
		if !r.Success {
			log.Printf("faile to get sentiment for %s\n", r.Text)
			num_failed += 1
		}
		req_time := r.Duration.Seconds() * 1000.0
		total_req_time_ms += req_time
		if req_time < min_req_time_ms {
			min_req_time_ms = req_time
		}
		if req_time > max_req_time_ms {
			max_req_time_ms = req_time
		}
	}
	elapsed := time.Since(t0)
	log.Printf("Did %d requests in %s (%g req/s)", *limit, elapsed, float64(*limit)/elapsed.Seconds())
	log.Printf("Req time: Avg %gms, Max %gms, Min %gms", total_req_time_ms, max_req_time_ms, min_req_time_ms)
}
