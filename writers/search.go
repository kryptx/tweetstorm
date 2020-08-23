package writers

import (
	"github.com/dghubble/go-twitter/twitter"
)

// IndexTweetWriter writes tweets to a search index
type IndexTweetWriter struct {
	Indexer SearchIndexer
}

// SearchIndexer is the interface that is consumed by the writer
type SearchIndexer interface {
	Index(tweet *twitter.Tweet) error
}

func (writer *IndexTweetWriter) Write(tweet *twitter.Tweet) <-chan error {
	result := make(chan error)
	go func() {
		result <- writer.Indexer.Index(tweet)
	}()
	return result
}
