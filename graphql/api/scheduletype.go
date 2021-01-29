package api

import "github.com/graphql-go/graphql"

var scheduleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Schedule",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"dueAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"url": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"method": &graphql.Field{
			Type: graphql.NewNonNull(httpMethodType),
		},
		"headers": &graphql.Field{
			Type: stringMapType,
		},
		"body": &graphql.Field{
			Type: graphql.String,
		},
	},
})
