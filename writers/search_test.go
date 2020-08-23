package writers_test

import (
	"errors"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/kryptx/tweetstorm/writers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIndexer struct {
	mock.Mock
}

func (indexer *MockIndexer) Index(tweet *twitter.Tweet) <-chan error {
	result := make(chan error)
	args := indexer.Called(tweet)
	go func() {
		result <- args.Error(0)
	}()
	return result
}

func setupSearchTest(err error) (writers.IndexTweetWriter, *twitter.Tweet) {
	indexer := new(MockIndexer)
	indexer.On("Index", mock.Anything).Return(err)
	writer := writers.IndexTweetWriter{Indexer: indexer}
	tweet := &twitter.Tweet{Text: "Foo"}
	return writer, tweet
}

func TestIndexWriter_Success_ReturnsNilError(t *testing.T) {
	writer, tweet := setupSearchTest(nil)
	err := <-writer.Write(tweet)
	assert.Nil(t, err)
}

func TestIndexWriter_Error_ReturnsError(t *testing.T) {
	errorMessage := "Mock error"
	writer, tweet := setupSearchTest(errors.New(errorMessage))
	err := <-writer.Write(tweet)
	assert.Equal(t, err.Error(), errorMessage)
}
