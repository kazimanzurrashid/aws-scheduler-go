package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
)

func (f *Factory) Create() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"dueAt": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"url": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"method": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(httpMethodType),
			},
			"headers": &graphql.ArgumentConfig{
				Type: stringMapType,
			},
			"body": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var input storage.CreateInput

			if err := loadStruct(p.Args, &input); err != nil {
				return nil, fmt.Errorf("invalid input")
			}

			if input.DueAt.Before(time.Now()) {
				return nil, fmt.Errorf("dueAt must be in future")
			}

			if input.URL == "" {
				return nil, fmt.Errorf("url is required")
			}

			if _, err := url.ParseRequestURI(input.URL); err != nil {
				return nil, fmt.Errorf("invalid url")
			}

			return f.storage.Create(p.Context, input)
		},
		Type: graphql.ID,
	}
}
