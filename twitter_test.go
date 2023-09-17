package gotweets_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	gotweets "github.com/RyoGreen/gotweet"
)

var c *gotweets.Client

func init() {
	c = gotweets.NewClient(os.Getenv("consumer_key"), os.Getenv("consumer_secret"), os.Getenv("access_token"), os.Getenv("access_token_secret"))
}

func TestUnauthorized(t *testing.T) {
	unauthorizedClient := gotweets.NewClient("", os.Getenv("consumer_secret"), os.Getenv("access_token"), os.Getenv("access_token_secret"))
	_, err := unauthorizedClient.Tweet(&gotweets.Options{
		Text: "test",
	})
	if err.Error() != "Unauthorized" {
		t.Errorf("Test failed, test name: Unauthorized Client, Expected error : %s, Actucal error %s", "Unauthorized", err.Error())
	} else {
		fmt.Println("Test passed")
	}
}

func TestTweets(t *testing.T) {
	tests := []struct {
		name          string
		option        *gotweets.Options
		expectedText  string
		expectedError error
	}{
		{
			name: "success tweet",
			option: &gotweets.Options{
				Text: "test",
			},
			expectedText:  "test",
			expectedError: nil,
		},
		{
			name: "error blank text",
			option: &gotweets.Options{
				Text: "",
			},
			expectedText:  "",
			expectedError: fmt.Errorf("error: One or more Text parameters is missing"),
		},
		// Not permitted to create an exclusive Tweet.
		// {
		// 	name: "success tweet for super follower only",
		// 	option: &gotweets.Options{
		// 		Text:                  "test for super follower only",
		// 		ForSuperFollowersOnly: true,
		// 	},
		// 	expectedText:  "test for super follower only",
		// 	expectedError: nil,
		// },
		{
			name: "success tweet with super follower only option is false",
			option: &gotweets.Options{
				Text:                  "test super follower only is false",
				ForSuperFollowersOnly: false,
			},
			expectedText:  "test super follower only is false",
			expectedError: nil,
		},
		{
			name: "success pull tweet",
			option: &gotweets.Options{
				Text: "poll tweet test",
				Poll: &gotweets.Poll{
					DurationMinutes: 300,
					Options:         []string{"yes", "no", "maybe"},
				},
			},
			expectedText:  "poll tweet test",
			expectedError: nil,
		},
		{
			name: "error poll options is nil",
			option: &gotweets.Options{
				Text: "test poll options is nil",
				Poll: &gotweets.Poll{
					DurationMinutes: 300,
				},
			},
			expectedText:  "",
			expectedError: fmt.Errorf("error: One or more Poll parameters is missing"),
		},
		{
			name: "error poll duration minutes is 0",
			option: &gotweets.Options{
				Text: "poll duration minutes is 0",
				Poll: &gotweets.Poll{
					Options:         []string{"yes", "no", "maybe"},
					DurationMinutes: 0,
				},
			},
			expectedText:  "",
			expectedError: fmt.Errorf("error: One or more Poll parameters is missing"),
		},
		{
			name: "success tweet with geo place id",
			option: &gotweets.Options{
				Text: "geo tweet test",
				Geo:  &gotweets.Geo{PlaceID: "df51dec6f4ee2b2c"},
			},
			expectedText:  "geo tweet test",
			expectedError: nil,
		},
		{
			name: "error geo place id is blank",
			option: &gotweets.Options{
				Text: "test geo place id is blank",
				Geo:  &gotweets.Geo{},
			},
			expectedText:  "",
			expectedError: fmt.Errorf("error: One or more GEO parameters is missing"),
		},
		{
			name: "success quote tweet",
			option: &gotweets.Options{
				Text:         "quote tweet https://t.co/OhT7X42IW1",
				QuoteTweetId: "1455953449422516226",
			},
			expectedText:  "quote tweet",
			expectedError: nil,
		},
		{
			name: "success reply tweet",
			option: &gotweets.Options{
				Text: "@XDevelopers reply tweet",
				Reply: &gotweets.Reply{
					InReplyToTweetId:    "1455953449422516226",
					ExcludeReplyUserIds: []string{"6253282"},
				},
			},
			expectedText:  "reply tweet",
			expectedError: nil,
		},
		{
			name: "error reply tweet id is blank",
			option: &gotweets.Options{
				Text: "reply tweet with blank reply tweet id",
				Reply: &gotweets.Reply{
					InReplyToTweetId:    "",
					ExcludeReplyUserIds: []string{"6253282"},
				},
			},
			expectedText:  "reply tweet",
			expectedError: fmt.Errorf("error: One or more Reply parameters is missing"),
		},
		{
			name:         "error set media, poll, reply options",
			expectedText: "",
			option: &gotweets.Options{
				Text: "test",
				Reply: &gotweets.Reply{
					InReplyToTweetId:    "1455953449422516226",
					ExcludeReplyUserIds: []string{"6253282"},
				},
				Poll: &gotweets.Poll{
					DurationMinutes: 300,
					Options:         []string{"yes", "no", "maybe"},
				},
				Media: &gotweets.Media{
					MediaIds: []string{"931270813906374700"},
				},
			},
			expectedError: fmt.Errorf("error: Media and Quote Tweet ID, Poll are mutually exclusive"),
		},
		{
			name: "error set media, reply options",
			option: &gotweets.Options{
				Text: "test",
				Reply: &gotweets.Reply{
					InReplyToTweetId:    "1455953449422516226",
					ExcludeReplyUserIds: []string{"6253282"},
				},
				Media: &gotweets.Media{
					MediaIds: []string{"931270813906374700"},
				},
			},
			expectedText:  "",
			expectedError: fmt.Errorf("error: Media and Quote Tweet ID, Poll are mutually exclusive"),
		},
		{
			name: "error set poll, reply options",
			option: &gotweets.Options{
				Text: "test",
				Reply: &gotweets.Reply{
					InReplyToTweetId:    "1455953449422516226",
					ExcludeReplyUserIds: []string{"6253282"},
				},
				Poll: &gotweets.Poll{
					DurationMinutes: 300,
					Options:         []string{"yes", "no", "maybe"},
				},
			},
			expectedText:  "",
			expectedError: fmt.Errorf("error: Media and Quote Tweet ID, Poll are mutually exclusive"),
		},
		{
			name: "error set media, poll options",
			option: &gotweets.Options{
				Text: "test",
				Poll: &gotweets.Poll{
					DurationMinutes: 300,
					Options:         []string{"yes", "no", "maybe"},
				},
				Media: &gotweets.Media{
					MediaIds: []string{"931270813906374700"},
				},
			},
			expectedText:  "",
			expectedError: fmt.Errorf("error: Media and Quote Tweet ID, Poll are mutually exclusive"),
		},
	}
	for _, test := range tests {
		resp, err := c.Tweet(test.option)
		time.Sleep(time.Second * 1)
		fmt.Println(test.name)
		if err != nil {
			if err.Error() == test.expectedError.Error() {
				fmt.Println("Test passed")
			} else {
				t.Errorf("Test failed, test name: %s, Expected error: %v, Actual error: %v", test.name, test.expectedError, err)
			}
		} else {
			if resp.Text != test.expectedText {
				t.Errorf("Test failed, test name: %s, Expected text: %s, Actual text: %s", test.name, test.expectedText, resp.Text)
			} else {
				fmt.Println("Test passed")
			}
		}
	}
}
