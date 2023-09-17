# go-twitter-tweets-api
Twitter Tweeting Library for Go

## Instalation:
```
go get github.com/RyoGreen/gotweet
```

## Prerequisites
Before using this library, make sure you have the following:
- Twitter Developer Account
- Twitter API Key and Access Tokens

## Example

```
package main

import (
  "log"
  "fmt"
  
  gotweets "github.com/RyoGreen/gotweet"
)

func main() {
	c = gotweets.NewClient("consumer_key", "consumer_secret", "access_token", "access_token_secret")
	resp, err := c.Tweet(&gotweets.Options{
		Text: "sample text",
	})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp.Text)
}
```

## Tweet Options
you can configure various tweet options.
- **Direct messageDeep link** (string): Set a direct message deep link.
- **For super followers only** (bool): Specify whether the tweet should be visible only to super followers.
- **Geo** (geo struct): Provide geographical information (place ID)
- **Media** (media struct): Attach media to your tweet. You can specify media IDs and tagged user IDs.
- **Poll** (poll struct): Create a poll in your tweet. Set the duration in minutes and provide options.
- **Quote tweet** (string): Include a quote tweet by specifying the quote tweet's ID.
- **Reply** (reply struct): Configure reply settings, including exclusion of reply user IDs and the in-reply-to tweet ID.
- **Reply setting** (ReplaySettingType): Set the reply settings, including "mentionedUsers" and "following."

## License
MIT
