package gotweets

import (
	"bytes"
	"encoding/json"
	"errors"
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

var errApiRateLimit = errors.New("error: tweets api rate limit exceeded")

func (c *Client) Tweet(text string) (*Data, error) {
	options := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}
	optionStr, err := json.Marshal(options)
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
