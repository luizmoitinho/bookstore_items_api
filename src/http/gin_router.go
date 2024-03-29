package http_router

import (
	"github.com/gin-gonic/gin"
	"github.com/luizmoitinho/bookstore_items_api/src/server"
	"github.com/luizmoitinho/bookstore_items_api/src/services"
	"github.com/luizmoitinho/bookstore_oauth/errors"
)

type Client interface {
	Run() error
	Routes()
}

type router struct {
	host   string
	engine *gin.Engine
}

func NewClient(host string) Client {
	return &router{
		host:   host,
		engine: gin.New(),
	}
}

func (r *router) Run() error {
	if r.host == "" {
		return errors.NewError("host is required")
	}
	return r.engine.Run(r.host)
}

func (r *router) Routes() {
	itemHandler := server.NewItemHandler(services.NewItemService())

	r.engine.POST("item/", itemHandler.Create)
	r.engine.GET("item/:id", itemHandler.Get)
	r.engine.GET("search/", itemHandler.Search)

}
