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
		"status": &graphql.Field{
			Type: graphql.NewNonNull(scheduleStatusType),
		},
		"startedAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"completedAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"canceledAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"result": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
	},
})
