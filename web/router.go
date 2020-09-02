package web

import (
	"net/http"

	"github.com/kryptx/tweetstorm/web/status"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/mongo"
)

// HTTPResponder responds to HTTP requests
type HTTPResponder interface {
	Respond(w http.ResponseWriter, r *http.Request)
}

func registerRoute(route string, responder HTTPResponder) {
	http.HandleFunc(route, responder.Respond)
}

// HandleRequests assigns routes to route handlers
func HandleRequests(collection *mongo.Collection, client *elastic.Client) {
	registerRoute("/images", &TweetImageHTTPResponder{Collection: collection})
	registerRoute("/status", &status.HTTPResponder{
		Retrievers: map[string]status.Retriever{
			"elasticsearch": &status.ElasticsearchRetriever{
				Client: client,
			},
			"mongo": &status.MongoRetriever{
				MongoCollection: collection,
			},
		},
	})
	http.ListenAndServe(":3000", nil)
}
