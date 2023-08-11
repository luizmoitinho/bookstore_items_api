package services

import (
	"github.com/luizmoitinho/bookstore_items_api/src/domain/items"
	"github.com/luizmoitinho/bookstore_utils/rest_errors"
)

type Items interface {
	Create(items.Item) (*items.Item, *rest_errors.RestError)
	Get(string) (*items.Item, *rest_errors.RestError)
}

type service struct{}

func NewItemService() Items {
	return &service{}
}

func (s *service) Create(item items.Item) (*items.Item, *rest_errors.RestError) {
	return nil, rest_errors.NewBadRequestError("implement me")
}

func (s *service) Get(id string) (*items.Item, *rest_errors.RestError) {
	if id != "" {
		return &items.Item{Id: id}, nil
	}
	return nil, rest_errors.NewBadRequestError("implement me")
}
