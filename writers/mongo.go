package writers

import (
	"context"

	"github.com/dghubble/go-twitter/twitter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoTweetWriter writes the tweet to a mongo database
type MongoTweetWriter struct {
	Collection WritableMongoCollection
}

// WritableMongoCollection specifies the interface the tweet writer depends on
type WritableMongoCollection interface {
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

func (writer *MongoTweetWriter) Write(tweet *twitter.Tweet) error {
	_, err := writer.Collection.InsertOne(context.TODO(), tweet)
	return err
}
