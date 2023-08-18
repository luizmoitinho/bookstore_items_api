package app

import (
	client "github.com/luizmoitinho/bookstore_items_api/src/clients/elasticsearch"
	http_router "github.com/luizmoitinho/bookstore_items_api/src/http"
	"github.com/luizmoitinho/bookstore_items_api/src/logger"
)

func Start(host string) {
	//check ES connection

	if _, err := client.NewElasticSearch(); err != nil {
		logger.Error("error when trying connecto to elastic search: %v.\n", err)
		panic(err)
	}

	server := http_router.NewClient(host)
	server.Routes()
	if err := server.Run(); err != nil {
		logger.Error("error during starting service: %v.\n", err)
	}
}
