package api

import "github.com/graphql-go/graphql"

var dataRangeType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "DateRange",
	Fields: graphql.InputObjectConfigFieldMap{
		"from": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"to": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
	},
})
