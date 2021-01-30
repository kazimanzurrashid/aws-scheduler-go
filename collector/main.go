package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/kazimanzurrashid/aws-scheduler-go/collector/storage"
)

var database storage.Storage

func handler(ctx context.Context) error {
	return database.Update(ctx)
}

func init() {
	ses := session.Must(session.NewSession())

	ddbc := dynamodb.New(ses)
	xray.AWS(ddbc.Client)

	database = storage.NewDatabase(ddbc)
}

func main() {
	lambda.Start(handler)
}
