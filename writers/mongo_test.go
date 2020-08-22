package writers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/kryptx/tweetstorm/writers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct {
	mock.Mock
}

func (collection *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := collection.Called(ctx, document, opts)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func TestMongoWriter_Success_ReturnsNilError(t *testing.T) {
	collection := new(MockCollection)
	collection.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.InsertOneResult{InsertedID: primitive.ObjectID{1}}, nil)
	writer := writers.MongoTweetWriter{Collection: collection}
	tweet := twitter.Tweet{Text: "Foo"}
	err := writer.Write(&tweet)
	assert.Nil(t, err)
}

func TestMongoWriter_Error_ReturnsError(t *testing.T) {
	errorMessage := "Mock error"
	collection := new(MockCollection)
	collection.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.InsertOneResult{InsertedID: nil}, errors.New(errorMessage))
	writer := writers.MongoTweetWriter{Collection: collection}
	tweet := twitter.Tweet{Text: "Foo"}
	err := writer.Write(&tweet)
	assert.Equal(t, err.Error(), errorMessage)
}
