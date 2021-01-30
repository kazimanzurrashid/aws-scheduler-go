package services

import "github.com/aws/aws-lambda-go/events"

type RequestInput struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    string
}

func CreateRequestInput(
	attributes map[string]events.DynamoDBAttributeValue) *RequestInput {

	input := RequestInput{
		URL:    attributes["url"].String(),
		Method: attributes["method"].String(),
	}

	if attr, found := attributes["headers"]; found && !attr.IsNull() {
		for k, v := range attr.Map() {
			input.Headers[k] = v.String()
		}
	} else {
		for k, v := range map[string]string{
			"Accept":       "application/json;charset=utf-8",
			"Content-Type": "application/json;charset=utf-8",
		} {
			input.Headers[k] = v
		}
	}

	if attr, found := attributes["body"]; found && !attr.IsNull() {
		input.Body = attr.String()
	}

	return &input
}
