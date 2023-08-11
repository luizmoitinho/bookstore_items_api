package main

import (
	"fmt"

	"github.com/luizmoitinho/bookstore_items_api/src/app"
	"github.com/luizmoitinho/bookstore_items_api/src/config"
)

func init() {
	config.Load(".env")
}

func main() {
	app.Start(fmt.Sprintf(":%d", config.Propertie.PORT))
}
