package gotweets_test

import (
	"os"
	"testing"

	gotweets "github.com/RyoGreen/gotweet"
)

var c *gotweets.Client

func init() {
	c = gotweets.NewClient(os.Getenv("consumer_key"), os.Getenv("consumer_secret"), os.Getenv("access_token"), os.Getenv("access_token_secret"))
}

func TestTweets(t *testing.T) {
	tests := []struct {
		name         string
		expectedText string
		option       *gotweets.Options
	}{
		{
			name:         "",
			expectedText: "test data",
			option:       &gotweets.Options{},
		},
	}
	for _, test := range tests {
		resp, err := c.Tweet(test.option)
		if err != nil {

		}
		if resp.Id == "" {
			t.Error(err)
		}
		if resp.RemainingRateLimit == 0 {
			t.Error(err)
		}
		if resp.Text != test.expectedText {
			t.Error(err)
		}
	}
}
