package writers

import (
	"context"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/olivere/elastic/v7"
)

// ElasticTweetWriter writes tweets to ElasticSearch
type ElasticTweetWriter struct {
	Client *elastic.Client
}

type jsonTweet struct {
	Username string    `json:"screenname"`
	Realname string    `json:"realname"`
	Text     string    `json:"message"`
	Image    string    `json:"image,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Hashtags []string  `json:"hashtags,omitempty"`
	Location string    `json:"location,omitempty"`
}

func (writer *ElasticTweetWriter) Write(tweet *twitter.Tweet) error {
	ctx := context.Background()
	time, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
	tags := []string{}
	for _, tag := range tweet.Entities.Hashtags {
		tags = append(tags, tag.Text)
	}

	tweetJSON := jsonTweet{
		Username: tweet.User.ScreenName,
		Realname: tweet.User.Name,
		Text:     tweet.Text,
		Created:  time,
		Hashtags: tags,
		Location: tweet.User.Location,
	}

	if len(tweet.Entities.Media) > 0 && tweet.Entities.Media[0].Type == "photo" {
		tweetJSON.Image = tweet.Entities.Media[0].MediaURL
	}

	_, err := writer.Client.Index().
		Index("twitter").
		Id(tweet.IDStr).
		BodyJson(tweetJSON).
		Do(ctx)

	return err
}
