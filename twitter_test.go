package gotweets_test

import (
	"os"
	"testing"

	gotweets "github.com/RyoGreen/gotweet"
)

func TestTweets(t *testing.T) {
	c := gotweets.NewClient(os.Getenv("consumer_key"), os.Getenv("consumer_secret"), os.Getenv("access_token"), os.Getenv("access_token_secret"))
	tests := []struct {
		text         string
		expectedText string
	}{
		{
			text:         "test data",
			expectedText: "test data",
		},
	}

	for _, v := range tests {
		respData, err := c.Tweet(v.text)
		if err != nil {
			t.Error(err)
			return
		}
		if respData.Text != v.expectedText {
			t.Errorf("Test failed. text: %v, Expected resutlt : %v, Actucal result %v", v.text, v.expectedText, respData.Text)
			return
		}
	}
}
