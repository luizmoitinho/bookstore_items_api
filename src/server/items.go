package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luizmoitinho/bookstore_items_api/src/domain/items"
	"github.com/luizmoitinho/bookstore_items_api/src/logger"
	"github.com/luizmoitinho/bookstore_items_api/src/services"
	"github.com/luizmoitinho/bookstore_oauth/errors"
	"github.com/luizmoitinho/bookstore_oauth/oauth"
)

type Items interface {
	Get(*gin.Context)
	Create(*gin.Context)
}

type itemHandler struct {
	service services.Items
}

func NewItemHandler(service services.Items) Items {
	return &itemHandler{
		service: service,
	}
}

func (handler *itemHandler) Create(c *gin.Context) {
	if err := oauth.Authenticate(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}

	accessToken, errCreate := handler.service.Create(items.Item{})
	if errCreate != nil {
		c.JSON(errCreate.Status, errCreate)
		return
	}
	c.JSON(http.StatusCreated, accessToken)
}

func (handler *itemHandler) Get(c *gin.Context) {
	accessToken, err := handler.service.Get(c.Param("id"))
	if err != nil {
		logger.Error("error during get item request", errors.NewError(err.Error))
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, accessToken)
}
