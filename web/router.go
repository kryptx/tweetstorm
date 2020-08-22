package web

import (
	"net/http"

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
	registerRoute("/status", &StatusHTTPResponder{
		ESClient:        client,
		MongoCollection: collection,
	})
	http.ListenAndServe(":3000", nil)
}
