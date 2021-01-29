package api

import (
	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
)

type Factory struct {
	storage storage.Storage
}

func NewFactory(storage storage.Storage) *Factory {
	return &Factory{storage}
}
