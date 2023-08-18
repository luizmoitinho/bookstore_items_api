package services

import (
	"fmt"

	client "github.com/luizmoitinho/bookstore_items_api/src/clients/elasticsearch"
	"github.com/luizmoitinho/bookstore_items_api/src/domain/items"
	"github.com/luizmoitinho/bookstore_items_api/src/logger"
	"github.com/luizmoitinho/bookstore_utils/rest_errors"
)

type Items interface {
	Create(items.Item) (*items.Item, *rest_errors.RestError)
	Get(string) (*items.Item, *rest_errors.RestError)
	Search(string) (*[]items.Item, *rest_errors.RestError)
}

type service struct{}

func NewItemService() Items {
	return &service{}
}

func (s *service) Create(item items.Item) (*items.Item, *rest_errors.RestError) {
	es, err := client.NewElasticSearch()
	if err != nil {
		logger.Error("error when trying connect with database", err)
		return nil, rest_errors.NewInternalServerError("error when trying connect with database", err)
	}
	itemDAO := items.NewItemDAO(es)
	if err := itemDAO.Save(item); err != nil {
		logger.Error("error when trying to save item", err)
		return nil, rest_errors.NewBadRequestError(fmt.Sprintf("error when trying to save item: %v", err))
	}

	return &item, nil
}

func (s *service) Get(id string) (*items.Item, *rest_errors.RestError) {
	es, err := client.NewElasticSearch()
	if err != nil {
		logger.Error("error when trying connect with database", err)
		return nil, rest_errors.NewInternalServerError("error when trying connect with database", err)
	}
	itemDAO := items.NewItemDAO(es)
	response, err := itemDAO.Get(id)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(fmt.Sprintf("error was receive during get by id: %v", err.Error()))
	}
	return response, nil
}

func (s *service) Search(desc string) (*[]items.Item, *rest_errors.RestError) {
	es, err := client.NewElasticSearch()
	if err != nil {
		logger.Error("error when trying connect with database", err)
		return nil, rest_errors.NewInternalServerError("error when trying connect with database", err)
	}
	itemDAO := items.NewItemDAO(es)

	response, err := itemDAO.Search(desc)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(fmt.Sprintf("error was receive during search: %v", err.Error()))
	}

	return response, nil
}
