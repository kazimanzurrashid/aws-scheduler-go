package api

import "github.com/graphql-go/graphql"

var scheduleListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ScheduleList",
	Fields: graphql.Fields{
		"schedule": &graphql.Field{
			Type: graphql.NewList(scheduleType),
		},
		"nextKey": &graphql.Field{
			Type: scheduleListNextKeyType,
		},
	},
})

