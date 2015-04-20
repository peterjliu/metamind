package metamind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"net/http"
)

const LanguageApi = "https://www.metamind.io/language/classify"

type textEndpoint struct {
	ClassifierId int    `json:"classifier_id"`
	Value        string `json:"value"`
}

type SentiClient struct {
	client *http.Client
	apikey string
}

type PredictedClass struct {
	ClassId   int     `json:"class_id"`
	ClassName string  `json:"class_name"`
	Prob      float32 `json:"prob"`
}

type SentiResp struct {
	Predictions []PredictedClass `json:"predictions"`
}

func NewSentimentClient(apikey string) *SentiClient {
	s := &SentiClient{
		client: &http.Client{},
		apikey: apikey,
	}
	return s
}

func (s *SentiClient) Classify(text string) (*SentiResp, error) {
	// 155 is sentiment classifier in demo
	testText := textEndpoint{ClassifierId: 155, Value: text}
	b, err := json.Marshal(testText)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", LanguageApi, bytes.NewReader(b))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", s.apikey))
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	return parseResp(resp)
}

func parseBody(body []byte) (*SentiResp, error) {
	var parsedResp SentiResp
	err := json.Unmarshal(body, &parsedResp)
	if err != nil {
		log.Print("error unmarshalling")
		return &parsedResp, err
	}
	return &parsedResp, nil
}

func parseResp(resp *http.Response) (*SentiResp, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseBody(body)
}

type SentimentResp struct {
	Text       string
	Prediction *SentiResp
	Success    bool
	Duration   time.Duration
}

// Get sentiment using metamind client for a string and put it in a channel.
func GetSentiment(client *SentiClient, text string, sentiments chan *SentimentResp) error {
	t0 := time.Now()
	senti, err := client.Classify(text)
	if err != nil {
		return err
	}
	sentiments <- &SentimentResp{Text: text, Prediction: senti, Success: err == nil, Duration: time.Since(t0)}
	return nil
}

func PrintSentiments(testText string, senti *SentiResp) {
	fmt.Printf("Sentiment for: %s\n", testText)
	for _, p := range senti.Predictions {
		fmt.Printf("%s %g\n", p.ClassName, p.Prob)
	}
	fmt.Printf("\n")
}
