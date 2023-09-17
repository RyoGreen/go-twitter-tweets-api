package gotweets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dghubble/oauth1"
)

type Client struct {
	c http.Client
}

func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *Client {
	c := oauth1.NewConfig(consumerKey, consumerSecret)
	t := oauth1.NewToken(accessToken, accessTokenSecret)
	return &Client{*c.Client(oauth1.NoContext, t)}
}

type Data struct {
	Id                 string `json:"id"`
	Text               string `json:"text"`
	RemainingRateLimit int
}

type Response struct {
	Data   Data   `json:"data"`
	Detail string `json:"detail"`
}

type ReplaySettingType string

const (
	MentionedUsers     ReplaySettingType = "mentionedUsers"
	MentionedFollowing ReplaySettingType = "following"
)

type Options struct {
	Text                  string            `json:"text"`
	DirectMessageDeepLink string            `json:"direct_message_deep_link,omitempty"`
	ForSuperFollowersOnly bool              `json:"for_super_followers_only,omitempty"`
	Geo                   *Geo              `json:"geo,omitempty"`
	Media                 *Media            `json:"media,omitempty"`
	Poll                  *Poll             `json:"poll,omitempty"`
	QuoteTweetId          string            `json:"quote_tweet_id,omitempty"`
	Reply                 *Reply            `json:"reply,omitempty"`
	ReplySettings         ReplaySettingType `json:"reply_settings,omitempty"`
}

type Geo struct {
	PlaceID string `json:"place_id"`
}

type Media struct {
	MediaIds      []string `json:"media_ids"`
	TaggedUserIds []string `json:"tagged_user_ids,omitempty"`
}

type Poll struct {
	DurationMinutes uint     `json:"duration_minutes"`
	Options         []string `json:"options"`
}

type Reply struct {
	ExcludeReplyUserIds []string `json:"exclude_reply_user_ids,omitempty"`
	InReplyToTweetId    string   `json:"in_reply_to_tweet_id"`
}

var errApiRateLimit = errors.New("error: tweets api rate limit exceeded")

const errOptionStr = "error: One or more %v parameters is missing"

func (c *Client) Tweet(o *Options) (*Data, error) {
	if o.Text == "" {
		return nil, fmt.Errorf(errOptionStr, "Text")
	}
	if o.Geo != nil && o.Geo.PlaceID == "" {
		return nil, fmt.Errorf(errOptionStr, "GEO")
	}
	if o.Media != nil && o.Media.MediaIds == nil {
		return nil, fmt.Errorf(errOptionStr, "Media")
	}
	if o.Poll != nil && (o.Poll.DurationMinutes == 0 || o.Poll.Options == nil) {
		return nil, fmt.Errorf(errOptionStr, "Poll")
	}
	if o.Reply != nil && o.Reply.InReplyToTweetId == "" {
		return nil, fmt.Errorf(errOptionStr, "Reply")
	}
	if (o.Media != nil && (o.Poll != nil || o.Reply != nil)) || (o.Poll != nil && (o.Media != nil || o.Reply != nil)) {
		return nil, errors.New("error: Media and Quote Tweet ID, Poll are mutually exclusive")
	}

	optionStr, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://api.twitter.com/2/tweets", bytes.NewBuffer(optionStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var respData Response
	if err := json.NewDecoder(res.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if respData.Detail != "" {
		return nil, errors.New(respData.Detail)
	}
	remainingRateLimit, err := strconv.Atoi(res.Header.Get("X-App-Limit-24hour-Remaining"))
	if err != nil {
		return nil, err
	}
	if remainingRateLimit <= 0 {
		return nil, errApiRateLimit
	}
	respData.Data.RemainingRateLimit = remainingRateLimit
	return &respData.Data, nil
}
