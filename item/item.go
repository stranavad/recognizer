package item

import "recognizer/types"

type Service struct {
	types.ServiceConfig
}

func NewItemService(config types.ServiceConfig) Service {
	return Service{config}
}
