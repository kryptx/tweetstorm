package status

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoStatus represents a status of mongodb
type mongoStatus struct {
	Ok        bool  `json:"ok"`
	Documents int64 `json:"documents"`
}

// MongoRetriever implements StatusRetriever for Mongo
type MongoRetriever struct {
	MongoCollection *mongo.Collection
}

// Retrieve returns the status in a channel
func (r *MongoRetriever) Retrieve(out chan<- interface{}) {
	status := mongoStatus{}
	documents, err := r.MongoCollection.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		log.Output(0, err.Error())
		out <- status
		return
	}

	status.Ok = true
	status.Documents = documents
	out <- status
}
