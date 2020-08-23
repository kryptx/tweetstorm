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
	GetJSONIndexer func(body interface{}, itemId string) ElasticsearchJSONIndexer
}

// ElasticsearchJSONIndexer indexes a JSON object in Elasticsearch
// it is implemented by the elastic client's IndexService type
type ElasticsearchJSONIndexer interface {
	Do(context.Context) (*elastic.IndexResponse, error)
}

// Index adds the tweet to the search index
func (indexer *ElasticsearchTweetIndexer) Index(tweet *twitter.Tweet) <-chan error {
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

	result := make(chan error)
	go func() {
		_, err := indexer.GetJSONIndexer(tweetJSON, tweet.IDStr).Do(context.Background())
		result <- err
	}()

	return result
}
