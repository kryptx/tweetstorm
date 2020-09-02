package status

import (
	"context"
	"log"

	"github.com/olivere/elastic/v7"
)

// elasticsearchStatus represents a status of the elasticsearch cluster
type elasticsearchStatus struct {
	Status       string `json:"status"`
	ActiveShards int    `json:"active_shards"`
}

// ElasticsearchRetriever implements StatusRetriever for Elasticsearch
type ElasticsearchRetriever struct {
	Client *elastic.Client
}

// Retrieve returns the status in a channel
func (r *ElasticsearchRetriever) Retrieve(out chan<- interface{}) {
	ctx := context.Background()
	result := elasticsearchStatus{}
	res, err := r.Client.
		ClusterHealth().
		Do(ctx)

	if err != nil {
		log.Output(0, err.Error())
		out <- result
	}

	result.Status = res.Status
	result.ActiveShards = res.ActiveShards

	out <- result
}
