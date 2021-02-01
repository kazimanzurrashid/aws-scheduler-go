package api

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

func (f *Factory) Cancel() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Args["id"].(string)

			if id == "" {
				return false, fmt.Errorf("id is required")
			}

			return f.storage.Cancel(p.Context, id)
		},
		Type: graphql.Boolean,
	}
}
