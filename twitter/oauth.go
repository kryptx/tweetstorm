package twitter

import (
	"context"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/kryptx/tweetstorm/config"
)

// GetClient returns an HTTP Client configured to authenticate with Twitter
func GetClient(config config.TwitterAuthConfig) *http.Client {
	ctx := context.Background()
	oauthConfig := oauth1.NewConfig(config.APIKey, config.APISecret)
	token := oauth1.NewToken(config.AccessToken, config.AccessTokenSecret)
	return oauthConfig.Client(ctx, token)
}
