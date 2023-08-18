package items

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	client "github.com/luizmoitinho/bookstore_items_api/src/clients/elasticsearch"
	"github.com/luizmoitinho/bookstore_items_api/src/logger"
)

const (
	indexItems = "items"
)

type ItemDAO struct {
	es client.ElasticSearch
}

type parserESReponse struct {
	Hits parserESReponseHits `json:"hits"`
}

type parserESReponseHits struct {
	Hits []esReponseHitsSource `json:"hits"`
}

type esReponseHitsSource struct {
	Source Item `json:"_source"`
}

func NewItemDAO(es client.ElasticSearch) ItemDAO {
	return ItemDAO{
		es: es,
	}
}

func (i *ItemDAO) Save(item Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		logger.Error("error during marshal in save item dao:", err)
		return err
	}
	i.es.Index(indexItems, data)

	return nil
}

func (i *ItemDAO) Get(id string) (*Item, error) {
	if id == "" {
		logger.Error("id is required", nil)
		return nil, errors.New("id is required")
	}
	response, err := i.es.Get(indexItems, id)
	if err != nil {
		return nil, err
	}

	items, err := parserItems(response.Body)
	if err != nil {
		return nil, err
	}

	if len(*items) == 1 {
		return &(*items)[0], nil
	}
	return nil, fmt.Errorf("was unexpected items was received: %v", *items)
}

func (i *ItemDAO) Search(desc string) (*[]Item, error) {
	if desc == "" {
		logger.Error("desc is required", nil)
		return nil, errors.New("desc is required")
	}
	response, err := i.es.Search(indexItems, desc)
	if err != nil {
		return nil, err
	}

	items, err := parserItems(response.Body)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func parserItems(body io.ReadCloser) (*[]Item, error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		logger.Error("error during read all elastic search response", err)
		return nil, err
	}

	parser := parserESReponse{}
	if err := json.Unmarshal(bytes, &parser); err != nil {
		logger.Error("error during unsmarshall elastic search response to map", err)
		return nil, err
	}

	_len := len(parser.Hits.Hits)
	response := make([]Item, _len)
	for index, p := range parser.Hits.Hits {
		response[index] = p.Source
	}
	return &response, nil
}
