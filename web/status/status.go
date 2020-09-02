package status

import (
	"net/http"

	"github.com/kryptx/tweetstorm/web/json"
)

/******************************
 * STATUS HTTP RESPONDER
 **/

// Retriever is an interface for asynchronous retrieval of status data from an external service or app
type Retriever interface {
	Retrieve(chan<- interface{})
}

// HTTPResponder sends HTTP responses containing the app's status
type HTTPResponder struct {
	json.Writer
	Retrievers map[string]Retriever
}

// Respond sends an HTTP response the app's status in JSON (implements web.HTTPResponder)
func (writer *HTTPResponder) Respond(w http.ResponseWriter, r *http.Request) {
	channels := map[string]chan interface{}{}
	results := map[string]interface{}{}
	for name, r := range writer.Retrievers {
		channels[name] = make(chan interface{})
		go r.Retrieve(channels[name])
	}

	for name, s := range channels {
		results[name] = <-s
	}

	writer.WriteJSON(w, 200, results)
}
