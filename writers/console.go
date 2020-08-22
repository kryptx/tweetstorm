package writers

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

// ConsoleTweetWriter is a TweetWriter that writes a tweet's text to the console
type ConsoleTweetWriter struct{}

func (writer *ConsoleTweetWriter) Write(tweet *twitter.Tweet) error {
	_, err := fmt.Println(tweet.Text)
	return err
}
