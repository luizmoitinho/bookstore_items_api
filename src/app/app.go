package app

import (
	"encoding/json"
	"fmt"

	client "github.com/luizmoitinho/bookstore_items_api/src/clients/elasticsearch"
	"github.com/luizmoitinho/bookstore_items_api/src/domain/items"
	http_router "github.com/luizmoitinho/bookstore_items_api/src/http"
	"github.com/luizmoitinho/bookstore_items_api/src/logger"
)

func Start(host string) {

	elasticSearch, err := client.NewElasticSearch()
	if err != nil {
		logger.Error("error during starting api", err)
		panic(err)
	}

	item := items.Item{
		Id:                "123",
		Seller:            1,
		Title:             "Teste",
		Description:       items.Description{},
		Pictures:          []items.Picture{},
		Video:             "123",
		Price:             129,
		AvailableQuantity: 0,
		SoldQuantity:      0,
		Status:            "active",
	}

	data, _ := json.Marshal(item)
	if resp, err := elasticSearch.Index("teste", data); err != nil {
		logger.Info(fmt.Sprintf("response: %v | error: %v", resp, err))
	}

	server := http_router.NewClient(host)
	server.Routes()
	if err := server.Run(); err != nil {
		logger.Error("error during starting service: %v.\n", err)
	}

}
