package main

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	var uis []*services.UpdateInput
	var wg sync.WaitGroup

	for _, record := range e.Records {
		if record.EventName != "MODIFY" {
			continue
		}

		if status := record.Change.NewImage["status"];
			status.String() != services.ScheduleStatusQueued {
			continue
		}

		wg.Add(1)
		go func(attrs map[string]events.DynamoDBAttributeValue) {
			defer wg.Done()

			ri := services.CreateRequestInput(attrs)
			ui := services.CreateUpdateInput(attrs)

			ui.StartedAt = time.Now().Unix()

			ro := httpClient.Request(ctx, ri)

			ui.Status = ro.Status
			ui.Result = ro.Result
			ui.CompletedAt = time.Now().Unix()

			uis = append(uis, ui)
		}(record.Change.NewImage)
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
