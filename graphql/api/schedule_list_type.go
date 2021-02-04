package api

import "github.com/graphql-go/graphql"

var scheduleListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ScheduleList",
	Fields: graphql.Fields{
		"schedules": &graphql.Field{
			Type: graphql.NewList(scheduleType),
		},
		"nextKey": &graphql.Field{
			Type: scheduleListNextKeyType,
		},
	},
})
