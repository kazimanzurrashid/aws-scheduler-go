package main

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/kazimanzurrashid/aws-scheduler-go/worker/services"
)

var (
	httpClient services.Client
	database   services.Storage
)

func handler(ctx context.Context, e events.DynamoDBEvent) error {
	var queuedRecords []events.DynamoDBEventRecord

	for _, record := range e.Records {
		if record.EventName != "MODIFY" {
			continue
		}

		if status := record.Change.NewImage["status"]; status.String() != services.ScheduleStatusQueued {
			continue
		}

		queuedRecords = append(queuedRecords, record)
	}

	if len(queuedRecords) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	uis := make([]*services.UpdateInput, len(queuedRecords))

	for i, record := range queuedRecords {
		wg.Add(1)

		go func(attrs map[string]events.DynamoDBAttributeValue, index int) {
			defer wg.Done()

			ri := services.CreateRequestInput(attrs)
			ui := services.CreateUpdateInput(attrs)

			ui.StartedAt = aws.Int64(time.Now().Unix())

			ro := httpClient.Request(ctx, ri)

			ui.Status = ro.Status
			ui.Result = aws.String(ro.Result)
			ui.CompletedAt = aws.Int64(time.Now().Unix())

			uis[index] = ui
		}(record.Change.NewImage, i)
	}

	wg.Wait()

	return database.Update(ctx, uis)
}

func init() {
	ses := session.Must(session.NewSession())

	ddbc := dynamodb.New(ses)
	xray.AWS(ddbc.Client)

	database = services.NewDatabase(ddbc)
	httpClient = services.NewHttpClient(
		xray.Client(retryablehttp.NewClient().StandardClient()))
}

func main() {
	lambda.Start(handler)
}
