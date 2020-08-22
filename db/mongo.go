package db

import (
	"context"
	"log"

	"github.com/kryptx/tweetstorm/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMongoCollection returns a connection to a mongo collection
func GetMongoCollection(config config.MongoConfig) *mongo.Collection {
	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(config.Database).Collection(config.Collection)
}
