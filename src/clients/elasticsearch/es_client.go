package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/luizmoitinho/bookstore_items_api/src/logger"
)

type ElasticSearch interface {
	AddIndex(index string) error
	Index(string, []byte) (*esapi.Response, error)
	Get(index, id string) (*esapi.Response, error)
	Search(index, id string) (*esapi.Response, error)
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

	if response, err := es.client.Cluster.Health(); err != nil {
		logger.Error("error when trying do a health report", err)
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("mismatch status code from elastic search, got %v and was expected %v", response.StatusCode, http.StatusOK)
		logger.Error("error during elastic search healtch check", err)
		return nil, err
	}

	return &es, nil
}

func (es *elasticSearch) Get(index string, id string) (*esapi.Response, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": id,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	return es.client.Search(es.client.Search.WithIndex(index), es.client.Search.WithBody(&buf), es.client.Search.WithFilterPath("hits.hits", "hits.total"))
}

func (es *elasticSearch) Search(index string, desc string) (*esapi.Response, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"title": desc,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	return es.client.Search(es.client.Search.WithIndex(index), es.client.Search.WithBody(&buf), es.client.Search.WithFilterPath("hits.hits", "hits.total"))
}

func (es *elasticSearch) Index(index string, value []byte) (*esapi.Response, error) {
	return es.client.Index(index, bytes.NewReader(value))
}

func (es *elasticSearch) AddIndex(index string) error {
	var indexes map[string]interface{}
	esapiResponse, err := es.client.Indices.Get([]string{index})
	if err != nil {
		return err
	}

	body, err := io.ReadAll(esapiResponse.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &indexes); err != nil {
		return err
	}

	if _, ok := indexes[index]; !ok {
		_, err := es.client.Indices.Create(index)
		return err
	}

	return nil
}
