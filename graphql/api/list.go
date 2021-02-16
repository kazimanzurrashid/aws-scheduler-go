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
				Type:         graphql.Int,
				DefaultValue: 25,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var input storage.ListInput

			if err := loadStruct(p.Args, &input); err != nil {
				return nil, fmt.Errorf("invalid input")
			}

			if input.DueAt != nil {
				if !input.DueAt.To.After(input.DueAt.From) {
					return nil, fmt.Errorf("dueAt to must be after dueAt from")
				}
			}

			if input.Limit < 1 || input.Limit > 100 {
				return nil, fmt.Errorf("limit must be between 1-100")
			}

			return f.storage.List(p.Context, input)
		},
		Type: scheduleListType,
	}
}
