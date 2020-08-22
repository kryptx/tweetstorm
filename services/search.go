package services

import (
	"context"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/olivere/elastic/v7"
)

type jsonTweet struct {
	Username string    `json:"screenname"`
	Realname string    `json:"realname"`
	Text     string    `json:"message"`
	Image    string    `json:"image,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Hashtags []string  `json:"hashtags,omitempty"`
	Location string    `json:"location,omitempty"`
}

// ElasticsearchTweetIndexer adds a tweet to an elasticsearch index
type ElasticsearchTweetIndexer struct {
	Client    *elastic.Client
	IndexName string
}

// Index adds the tweet to the search index
func (indexer *ElasticsearchTweetIndexer) Index(tweet *twitter.Tweet) error {
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

	_, err := indexer.Client.Index().
		Index(indexer.IndexName).
		Id(tweet.IDStr).
		BodyJson(tweetJSON).
		Do(ctx)

	return err
}
