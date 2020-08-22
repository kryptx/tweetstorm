package config

import (
	"os"

	cfg "go.uber.org/config"
)

// AppConfig contains the top-level configuration for the app
type AppConfig struct {
	Twitter       TwitterConfig `yaml:"twitter"`
	Mongo         MongoConfig   `yaml:"mongodb"`
	ElasticSearch ElasticConfig `yaml:"elasticsearch"`
}

// TwitterConfig contains configuration specific to Twitter
type TwitterConfig struct {
	FilterTerms []string          `yaml:"filterTerms"`
	Auth        TwitterAuthConfig `yaml:"auth"`
}

// TwitterAuthConfig contains values for authenticating with Twitter
type TwitterAuthConfig struct {
	APIKey            string `yaml:"APIKey"`
	APISecret         string `yaml:"APISecret"`
	AccessToken       string `yaml:"accessToken"`
	AccessTokenSecret string `yaml:"accessTokenSecret"`
}

// MongoConfig contains mongoDB configuration
type MongoConfig struct {
	URI        string `yaml:"URI"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

// ElasticConfig contains ElasticSearch configuration
type ElasticConfig struct {
	URI   string `yaml:"URI"`
	Index string `yaml:"index"`
}

// Load loads the config from ./config.yml with env var expansion
func Load(filename string) *AppConfig {
	provider, err := cfg.NewYAML(cfg.File(filename), cfg.Expand(os.LookupEnv))
	if err != nil {
		panic(err)
	}

	var c AppConfig
	if err := provider.Get("config").Populate(&c); err != nil {
		panic(err)
	}
	return &c
}
