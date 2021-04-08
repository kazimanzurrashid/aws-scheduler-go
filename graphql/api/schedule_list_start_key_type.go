package api

import "github.com/graphql-go/graphql"

var scheduleListStartKeyType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ScheduleListStartKey",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"dueAt": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
		"status": &graphql.InputObjectFieldConfig{
			Type: scheduleStatusType,
		},
	},
})
