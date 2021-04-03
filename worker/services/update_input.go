package services

import "github.com/aws/aws-lambda-go/events"

type UpdateInput struct {
	ID          string            `dynamodbav:"id"`
	DueAt       int64             `dynamodbav:"dueAt"`
	URL         string            `dynamodbav:"url"`
	Method      string            `dynamodbav:"method"`
	Headers     map[string]string `dynamodbav:"headers,omitempty"`
	Body        *string           `dynamodbav:"body,omitempty"`
	StartedAt   *int64            `dynamodbav:"startedAt"`
	CompletedAt *int64            `dynamodbav:"completedAt"`
	Status      string            `dynamodbav:"status"`
	Result      *string           `dynamodbav:"result"`
	CreatedAt   int64             `dynamodbav:"createdAt"`
}

func CreateUpdateInput(
	attributes map[string]events.DynamoDBAttributeValue) *UpdateInput {

	var body string

	if attr, found := attributes["body"]; found && !attr.IsNull() {
		body = attr.String()
	}

	headers := make(map[string]string)

	if attr, found := attributes["headers"]; found && !attr.IsNull() {
		for k, v := range attr.Map() {
			headers[k] = v.String()
		}
	}

	dueAt, _ := attributes["dueAt"].Integer()
	createdAt, _ := attributes["createdAt"].Integer()

	input := UpdateInput{
		ID:        attributes["id"].String(),
		DueAt:     dueAt,
		URL:       attributes["url"].String(),
		Method:    attributes["method"].String(),
		Headers:   headers,
		Body:      &body,
		CreatedAt: createdAt,
	}

	return &input
}
