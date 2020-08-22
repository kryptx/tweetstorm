package web

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// statusResult represents a result sent from the status endpoing
type statusResult struct {
	Elasticsearch elasticsearchStatus `json:"elasticsearch"`
	Mongo         mongoStatus         `json:"mongo"`
}

// elasticsearchStatus represents a status of the elasticsearch cluster
type elasticsearchStatus struct {
	Status       string `json:"status"`
	ActiveShards int    `json:"active_shards"`
}

// indexStatus represents some data about an index
type indexStatus struct {
	Status       string `json:"status"`
	ActiveShards int    `json:"active_shards"`
	Replicas     int    `json:"replicas"`
}

// mongoStatus represents a status of mongodb
type mongoStatus struct {
	Ok        bool  `json:"ok"`
	Documents int64 `json:"documents"`
}

// StatusHTTPResponder sends HTTP responses containing the app's status
type StatusHTTPResponder struct {
	MongoCollection *mongo.Collection
	ESClient        *elastic.Client
}

// Respond sends an HTTP response the app's status in JSON
func (writer *StatusHTTPResponder) Respond(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	esStatus := getElasticsearchStatus(writer.ESClient)
	mongoStatus := getMongoStatus(writer.MongoCollection)
	result := statusResult{
		Elasticsearch: esStatus,
		Mongo:         mongoStatus,
	}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(result)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(buf.Bytes())
}

func getElasticsearchStatus(client *elastic.Client) elasticsearchStatus {
	ctx := context.Background()
	result := elasticsearchStatus{}
	res, err := client.
		ClusterHealth().
		Do(ctx)

	if err != nil {
		log.Output(0, err.Error())
		return result
	}

	result.Status = res.Status
	result.ActiveShards = res.ActiveShards

	return result
}

func getMongoStatus(collection *mongo.Collection) mongoStatus {
	ctx := context.Background()
	result := mongoStatus{}
	documents, err := collection.CountDocuments(ctx, bson.M{})

	if err != nil {
		log.Output(0, err.Error())
		return result
	}

	result.Ok = true
	result.Documents = documents
	return result
}
