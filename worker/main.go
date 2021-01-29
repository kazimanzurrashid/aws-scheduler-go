package main

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/hashicorp/go-retryablehttp"

	"golang.org/x/sync/errgroup"
)

var httpClient *http.Client

func handle(ctx context.Context, record map[string]events.DynamoDBAttributeValue) error {
	var body string

	if attr, found := record["body"]; found && !attr.IsNull() {
		body = attr.String()
	}

	req, err := http.NewRequestWithContext(
		ctx,
		record["method"].String(),
		record["url"].String(),
		bytes.NewBufferString(body))

	if err != nil {
		return err
	}

	if attr, found := record["headers"]; found && !attr.IsNull() {
		for k, v := range attr.Map() {
			req.Header.Set(k, v.String())
		}
	} else {
		for k, v := range map[string]string{
			"Accept":       "application/json;charset=utf-8",
			"Content-Type": "application/json;charset=utf-8",
		} {
			req.Header.Set(k, v)
		}
	}

	if _, err = httpClient.Do(req); err != nil {
		return err
	}

	return nil
}

func handler(ctx context.Context, e events.DynamoDBEvent) error {
	g, _ := errgroup.WithContext(ctx)

	for _, record := range e.Records {
		if record.EventName != "REMOVE" {
			continue
		}

		if canceled, ok := record.Change.OldImage["canceled"]; ok && canceled.Boolean() {
			continue
		}

		g.Go(func() error {
			return handle(ctx, record.Change.OldImage)
		})
	}

	return g.Wait()
}

func init() {
	httpClient = xray.Client(retryablehttp.NewClient().StandardClient())
}

func main() {
	lambda.Start(handler)
}
