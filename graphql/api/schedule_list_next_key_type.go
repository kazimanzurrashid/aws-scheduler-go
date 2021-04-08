package api

import "github.com/graphql-go/graphql"

var scheduleListNextKeyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ScheduleListNextKey",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"dueAt": &graphql.Field{
			Type: graphql.Int,
		},
		"status": &graphql.Field{
			Type: scheduleStatusType,
		},
	},
})
