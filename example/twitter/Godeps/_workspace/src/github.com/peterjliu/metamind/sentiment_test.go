package metamind

import (
	"testing"
)

const testResp = `
{
  "predictions": [
    {
      "class_id": 2,
      "class_name": "positive",
      "prob": 0.992522780802063
    },
    {
      "class_id": 0,
      "class_name": "neutral",
      "prob": 0.0042648534368411265
    },
    {
      "class_id": 1,
      "class_name": "negative",
      "prob": 0.00321236576109588
    }
  ]
}
`

func TestJsonParse(t *testing.T) {
	sent, err := parseBody([]byte(testResp))
	if err != nil {
		t.Errorf("Couldn't parse testResp\n %s", err)
	}
	if len(sent.Predictions) != 3 {
		t.Errorf("Unexpected number of predictions")
	}
}
