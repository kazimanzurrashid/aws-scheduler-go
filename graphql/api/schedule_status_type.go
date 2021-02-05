package api

import (
	"github.com/graphql-go/graphql"

	"github.com/kazimanzurrashid/aws-scheduler-go/graphql/storage"
)

var scheduleStatusType = graphql.NewEnum(graphql.EnumConfig{
	Name: "ScheduleStatus",
	Values: map[string]*graphql.EnumValueConfig{
		"IDLE":      {Value: storage.ScheduleStatusIdle},
		"QUEUED":    {Value: storage.ScheduleStatusQueued},
		"SUCCEEDED": {Value: storage.ScheduleStatusSucceeded},
		"CANCELED":  {Value: storage.ScheduleStatusCanceled},
		"FAILED":    {Value: storage.ScheduleStatusFailed},
	},
})
