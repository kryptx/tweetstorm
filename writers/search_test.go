package writers_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/kryptx/tweetstorm/writers"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/**
 * Test high-level abstraction: mock indexer itself
 */

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

func setupHighLevelSearchTest(err error) (writers.IndexTweetWriter, *twitter.Tweet) {
	indexer := new(MockIndexer)
	indexer.On("Index", mock.Anything).Return(err)
	writer := writers.IndexTweetWriter{Indexer: indexer}
	tweet := &twitter.Tweet{Text: "Foo"}
	return writer, tweet
}

func TestIndexWriter_Success_ReturnsNilError(t *testing.T) {
	writer, tweet := setupHighLevelSearchTest(nil)
	err := <-writer.Write(tweet)
	assert.Nil(t, err)
}

func TestIndexWriter_Error_ReturnsError(t *testing.T) {
	errorMessage := "Mock error"
	writer, tweet := setupHighLevelSearchTest(errors.New(errorMessage))
	err := <-writer.Write(tweet)
	assert.Equal(t, err.Error(), errorMessage)
}

/**
 * Test TweetJSONIndexer.Index (mock JSONIndexer, which is essentially ElasticSearch)
 */

type MockJSONIndexerFactory struct {
	mock.Mock
	jsonIndexer writers.JSONIndexer
}

type MockJSONIndexer struct {
	mock.Mock
}

func (indexer *MockJSONIndexer) Do(ctx context.Context) (*elastic.IndexResponse, error) {
	args := indexer.Called(ctx)
	return args.Get(0).(*elastic.IndexResponse), args.Error(1)
}

func (indexer *MockJSONIndexerFactory) CreateJSONIndexer(body interface{}, itemId string) writers.JSONIndexer {
	args := indexer.Called(body, itemId)
	return args.Get(0).(writers.JSONIndexer)
}

func setupIndexerTest(response *elastic.IndexResponse, err error) (*writers.TweetJSONIndexer, *MockJSONIndexer, *twitter.Tweet) {
	jsonIndexer := new(MockJSONIndexer)
	jsonIndexer.On("Do", mock.Anything).Return(response, err)
	factory := &MockJSONIndexerFactory{jsonIndexer: jsonIndexer}
	factory.On("CreateJSONIndexer", mock.Anything, mock.Anything).Return(jsonIndexer)
	tweet := &twitter.Tweet{
		ID:    12345678,
		IDStr: "12345678",
		User: &twitter.User{
			ScreenName: "Burbage",
			Name:       "Jason Burbage",
			Location:   "Murica",
		},
		Text: "I just wrote a test!",
		Entities: &twitter.Entities{
			Hashtags: []twitter.HashtagEntity{},
		},
	}
	return &writers.TweetJSONIndexer{factory}, jsonIndexer, tweet
}

func TestIndex_Error_ReturnsErrInChannel(t *testing.T) {
	err := errors.New("test error")
	indexer, _, tweet := setupIndexerTest(nil, err)
	indexErr := <-indexer.Index(tweet)
	assert.Equal(t, err, indexErr)
}

func TestIndex_Success_PassesCorrectArguments(t *testing.T) {
	response := &elastic.IndexResponse{}
	indexer, jsonIndexer, tweet := setupIndexerTest(response, nil)

	_ = <-indexer.Index(tweet)
	assert.Equal(t, 1, len(jsonIndexer.Calls))

	mockFactory := indexer.Factory.(*MockJSONIndexerFactory)
	assert.Equal(t, 1, len(mockFactory.Calls))

	id := mockFactory.Calls[0].Arguments[1].(string)
	assert.Equal(t, "12345678", id)

	jsonObj := mockFactory.Calls[0].Arguments[0]
	json, jsonErr := json.Marshal(jsonObj)
	assert.Nil(t, jsonErr)
	assert.Contains(t, string(json), "\"I just wrote a test!\"")
}
