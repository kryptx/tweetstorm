package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/kryptx/tweetstorm/config"
	"github.com/olivere/elastic/v7"
)

// GetElasticClient returns an ElasticSearch client
func GetElasticClient(config config.ElasticConfig) *elastic.Client {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	client, err := elastic.NewClient(
		elastic.SetURL(config.URI),
		elastic.SetHealthcheckTimeoutStartup(60*time.Second),
		elastic.SetRetrier(
			elastic.NewBackoffRetrier(elastic.NewConstantBackoff(5*time.Second)),
		),
	)

	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(config.URI).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(config.URI)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(config.Index).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		mapping, err := ioutil.ReadFile(fmt.Sprintf("indexes/%v.json", config.Index))
		if err != nil {
			log.Fatal(err)
		}

		createIndex, err := client.CreateIndex(config.Index).Body(string(mapping)).Do(ctx)
		if err != nil {
			// Handle error
			log.Fatal(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return client
}
