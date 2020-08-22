package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TweetResult represents a result from mongodb about a tweet
type TweetResult struct {
	Image   TweetEntities `bson:"entities"`
	User    TweetUser     `bson:"user"`
	TweetID string        `bson:"idstr"`
}

// TweetUser contains data about a user that posted a tweet
type TweetUser struct {
	ScreenName string `bson:"screenname"`
}

// TweetEntities contains entities associated with a tweet
type TweetEntities struct {
	Media []TweetMedia `bson:"media"`
}

// TweetMedia contains data about media entities from a tweet
type TweetMedia struct {
	URL string `bson:"mediaurl"`
}

// TweetImageHTTPResponder sends HTTP responses containing tweets' images
type TweetImageHTTPResponder struct {
	Collection *mongo.Collection
}

// Respond sends an HTTP response containing images that link to the tweet
func (writer *TweetImageHTTPResponder) Respond(w http.ResponseWriter, r *http.Request) {
	var results []TweetResult
	ctx := context.Background()
	w.Header().Add("Content-type", "text/html")

	cursor, err := writer.Collection.Find(ctx,
		bson.M{
			"entities.media.type": "photo",
		})

	if err != nil {
		log.Output(0, err.Error())
	}

	if err = cursor.All(ctx, &results); err != nil {
		log.Output(0, err.Error())
	}

	for _, tweet := range results {
		w.Write([]byte(fmt.Sprintf("<a href=\"https://www.twitter.com/%v/status/%v\"><img src=\"%v\" /></a><br />\n", tweet.User.ScreenName, tweet.TweetID, tweet.Image.Media[0].URL)))
	}

}
