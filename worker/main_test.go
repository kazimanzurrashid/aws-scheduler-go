package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type fakeTransport struct {
	Requests []*http.Request
}

func (ft *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ft.Requests = append(ft.Requests, r)

	return &http.Response{}, nil
}

const (
	url    = "https://foo.bar/do"
	method = "POST"
	accept = "application/json;charset=utf-8"
	body   = "{ \"foo\": \"bar\" }"
)

var ctx = context.TODO()

func Test_handle_Success(t *testing.T) {
	ft := fakeTransport{}
	c := &http.Client{
		Transport: &ft,
	}

	httpClient = c

	err := handler(ctx, events.DynamoDBEvent{
		Records: []events.DynamoDBEventRecord{
			{
				EventName: "ADD",
			},
			{
				EventName: "REMOVE",
				Change: events.DynamoDBStreamRecord{
					OldImage: map[string]events.DynamoDBAttributeValue{
						"url":    events.NewStringAttribute(url),
						"method": events.NewStringAttribute(method),
						"headers": events.NewMapAttribute(map[string]events.DynamoDBAttributeValue{
							"accept": events.NewStringAttribute(accept),
						}),
						"body": events.NewStringAttribute(body),
					},
				},
			},
			{
				EventName: "REMOVE",
				Change: events.DynamoDBStreamRecord{
					OldImage: map[string]events.DynamoDBAttributeValue{
						"url":    events.NewStringAttribute(url),
						"method": events.NewStringAttribute(method),
					},
				},
			},
			{
				EventName: "REMOVE",
				Change: events.DynamoDBStreamRecord{
					OldImage: map[string]events.DynamoDBAttributeValue{
						"canceled": events.NewBooleanAttribute(true),
					},
				},
			},
		},
	})

	assert.Nil(t, err)

	for _, req := range ft.Requests {
		assert.Equal(t, url, req.URL.String())
		assert.Equal(t, method, req.Method)
		assert.Equal(t, accept, req.Header.Get("accept"))
	}
}
