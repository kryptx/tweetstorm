package writers

import (
	"context"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/olivere/elastic/v7"
)

const twitterTimeFormat = "Mon Jan 2 15:04:05 -0700 2006"

/**
 * High-level abstraction: writing tweets to a search index
 */

// IndexTweetWriter writes tweets to a search index
type IndexTweetWriter struct {
	Indexer TweetIndexer
}

func (writer *IndexTweetWriter) Write(tweet *twitter.Tweet) <-chan error {
	return writer.Indexer.Index(tweet)
}

/**
 * Mid-level abstraction: writing JSON to a search index
 */

type jsonTweet struct {
	Username string    `json:"screenname"`
	Realname string    `json:"realname"`
	Text     string    `json:"message"`
	Image    string    `json:"image,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Hashtags []string  `json:"hashtags,omitempty"`
	Location string    `json:"location,omitempty"`
}

// JSONIndexer indexes a JSON object in a search index
// it is implemented by the elastic client's IndexService type
type JSONIndexer interface {
	Do(context.Context) (*elastic.IndexResponse, error)
}

// JSONIndexerFactory creates indexers
type JSONIndexerFactory interface {
	CreateJSONIndexer(body interface{}, id string) JSONIndexer
}

// TweetIndexer indexes tweets in a search index
type TweetIndexer interface {
	Index(tweet *twitter.Tweet) <-chan error
}

// TweetJSONIndexer indexes tweets using JSON
type TweetJSONIndexer struct {
	Factory JSONIndexerFactory
}

// Index adds the tweet to the search index
func (indexer *TweetJSONIndexer) Index(tweet *twitter.Tweet) <-chan error {
	time, _ := time.Parse(twitterTimeFormat, tweet.CreatedAt)
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
		_, err := indexer.Factory.CreateJSONIndexer(tweetJSON, tweet.IDStr).Do(context.Background())
		result <- err
	}()

	return result
}

/**
 * Concrete implementation: elasticsearch-specific implementation
 */

// ElasticsearchJSONIndexerFactory creates indexers that index in elasticsearch
type ElasticsearchJSONIndexerFactory struct {
	Client    *elastic.Client
	IndexName string
}

// CreateJSONIndexer is the factory function for the ElasticsearchJSONIndexerFactory
func (factory *ElasticsearchJSONIndexerFactory) CreateJSONIndexer(body interface{}, id string) JSONIndexer {
	return factory.Client.Index().Index(factory.IndexName).Id(id).BodyJson(body)
}
