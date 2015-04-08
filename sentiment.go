package metamind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

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

func NewSentimentClient(apikey string) *SentiClient {
	s := &SentiClient{
		client: &http.Client{},
		apikey: apikey,
	}
	return s
}

func (s *SentiClient) Classify(text string) ([]byte, error) {
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
