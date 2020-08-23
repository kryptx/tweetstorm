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
	getJSONIndexer := func(body interface{}, id string) services.ElasticsearchJSONIndexer {
		return elasticClient.Index().Index(c.ElasticSearch.Index).Id(id).BodyJson(body)
	}
	tweetWriters := []twitter.TweetWriter{
		&writers.MongoTweetWriter{Collection: mongoCollection},
		&writers.IndexTweetWriter{
			Indexer: &services.ElasticsearchTweetIndexer{
				GetJSONIndexer: getJSONIndexer,
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
