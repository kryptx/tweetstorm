package main

import (
	"github.com/kryptx/tweetstorm/config"
	"github.com/kryptx/tweetstorm/db"
	"github.com/kryptx/tweetstorm/twitter"
	"github.com/kryptx/tweetstorm/web"
	writers "github.com/kryptx/tweetstorm/writers"
)

func main() {
	c := config.Load("config.yml")
	mongoCollection := db.GetMongoCollection(c.Mongo)
	elasticClient := db.GetElasticClient(c.ElasticSearch)
	indexer := &writers.TweetJSONIndexer{
		Factory: &writers.ElasticsearchJSONIndexerFactory{
			Client:    elasticClient,
			IndexName: c.ElasticSearch.Index,
		},
	}
	tweetWriters := []twitter.TweetWriter{
		&writers.MongoTweetWriter{Collection: mongoCollection},
		&writers.IndexTweetWriter{Indexer: indexer},
	}
	go web.HandleRequests(mongoCollection, elasticClient)
	twitter.StreamTweets(
		c.Twitter.FilterTerms,
		twitter.GetClient(c.Twitter.Auth),
		tweetWriters,
	)
}
