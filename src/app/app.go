package app

import (
	"log"

	http_router "github.com/luizmoitinho/bookstore_items_api/src/http"
)

func Start(host string) {
	server := http_router.NewClient(host)
	server.Routes()
	if err := server.Run(); err != nil {
		log.Fatalf("error during starting service: %v.\n", err)
	}

}
