package writers

import (
	"github.com/dghubble/go-twitter/twitter"
)

// IndexTweetWriter writes tweets to a search index
type IndexTweetWriter struct {
	Indexer TweetIndexer
}

// TweetIndexer is the interface that is consumed by the writer
type TweetIndexer interface {
	Index(tweet *twitter.Tweet) <-chan error
}

func (writer *IndexTweetWriter) Write(tweet *twitter.Tweet) <-chan error {
	return writer.Indexer.Index(tweet)
}
