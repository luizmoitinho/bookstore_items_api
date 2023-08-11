package client

import (
	"bytes"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/luizmoitinho/bookstore_items_api/src/logger"
)

type ElasticSearch interface {
	Index(string, []byte) (*esapi.Response, error)
}

type elasticSearch struct {
	client *elasticsearch.Client
}

func NewElasticSearch() (ElasticSearch, error) {
	var es elasticSearch

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200/"},
		Username:  "elastic",
		Password:  "123456",
	})
	if err != nil {
		logger.Error("error when trying create an elastic search client", err)
		return nil, err
	}

	es.client = client
	es.client.Cluster.Health.WithTimeout(10 * time.Second)
	es.client.Cluster.AllocationExplain.WithErrorTrace()
	es.client.Cluster.AllocationExplain.WithFilterPath()

	es.client.Indices.Create("items-api")

	return &es, nil
}

func (es *elasticSearch) Index(index string, value []byte) (*esapi.Response, error) {
	return es.client.Index(index, bytes.NewReader(value))
}
