package main

import (
	"github.com/kryptx/tweetstorm/config"
	"github.com/kryptx/tweetstorm/db"
	"github.com/kryptx/tweetstorm/services"
	"github.com/kryptx/tweetstorm/twitter"
	"github.com/kryptx/tweetstorm/web"
	writers "github.com/kryptx/tweetstorm/writers"
)

func main() {
	c := config.Load("config.yml")
	mongoCollection := db.GetMongoCollection(c.Mongo)
	elasticClient := db.GetElasticClient(c.ElasticSearch)
	tweetWriters := []twitter.TweetWriter{
		&writers.MongoTweetWriter{Collection: mongoCollection},
		&writers.IndexTweetWriter{
			Indexer: &services.ElasticsearchTweetIndexer{
				Client: elasticClient,
			},
		},
	}
	go web.HandleRequests(mongoCollection, elasticClient)
	twitter.StreamTweets(
		c.Twitter.FilterTerms,
		twitter.GetClient(c.Twitter.Auth),
		tweetWriters,
	)
}
