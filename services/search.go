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

	return indexJSONObject(tweetJSON, indexer.IndexName, tweet.IDStr, indexer.Client)
}

func indexJSONObject(obj interface{}, index string, id string, client *elastic.Client) error {
	_, err := client.Index().
		Index(index).
		Id(id).
		BodyJson(obj).
		Do(context.Background())

	return err
}
