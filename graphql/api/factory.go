package api

import (
	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
)

type Factory struct {
	storage storage.Storage
}

func NewFactory(storage storage.Storage) *Factory {
	return &Factory{storage}
}

func (f *Factory) Schema() (graphql.Schema, error) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Queries",
		Fields: graphql.Fields{
			"get": f.Get(),
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutations",
		Fields: graphql.Fields{
			"create": f.Create(),
			"cancel": f.Cancel(),
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
}
