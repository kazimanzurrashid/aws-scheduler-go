package storage

import "time"

type CreateInput struct {
	DueAt   time.Time         `json:"dueAt" dynamodbav:"dueAt,unixtime"`
	URL     string            `json:"url" dynamodbav:"url"`
	Method  string            `json:"method" dynamodbav:"method"`
	Headers map[string]string `json:"headers,omitempty" dynamodbav:"headers,omitempty"`
	Body    string            `json:"body,omitempty" dynamodbav:"body,omitempty"`
}
