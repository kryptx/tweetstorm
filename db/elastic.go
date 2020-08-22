package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kryptx/tweetstorm/config"
	"github.com/olivere/elastic/v7"
)

// GetElasticClient returns an ElasticSearch client
func GetElasticClient(config config.ElasticConfig) *elastic.Client {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient(elastic.SetURL(config.URI))
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
