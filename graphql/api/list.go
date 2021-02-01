package api

import (
	"fmt"
	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
)

func (f *Factory) List() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"status": &graphql.ArgumentConfig{
				Type: scheduleStatusType,
			},
			"dueAt": &graphql.ArgumentConfig{
				Type: dataRangeType,
			},
			"startKey": &graphql.ArgumentConfig{
				Type: scheduleListStartKeyType,
			},
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
				DefaultValue: 25,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var input storage.ListInput

			if err := loadStruct(p.Args, &input); err != nil {
				return nil, fmt.Errorf("invalid input")
			}

			return f.storage.List(p.Context, input)
		},
		Type: scheduleListType,
	}
}
