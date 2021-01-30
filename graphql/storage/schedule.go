package storage

import "time"

type Schedule struct {
	ID          string            `json:"id"`
	DueAt       time.Time         `json:"dueAt" dynamodbav:"dueAt,unixtime"`
	URL         string            `json:"url" dynamodbav:"url"`
	Method      string            `json:"method" dynamodbav:"method"`
	Headers     map[string]string `json:"headers,omitempty" dynamodbav:"headers,omitempty"`
	Body        string            `json:"body,omitempty" dynamodbav:"body,omitempty"`
	Status      string            `json:"status" dynamodbav:"status"`
	StartedAt   time.Time         `json:"startedAt,omitempty" dynamodbav:"startedAt,unixtime,omitempty"`
	CompletedAt time.Time         `json:"completedAt,omitempty" dynamodbav:"completedAt,unixtime,omitempty"`
	CanceledAt  time.Time         `json:"canceledAt,omitempty" dynamodbav:"canceledAt,unixtime,omitempty"`
	Result      string            `json:"result,omitempty" dynamodbav:"result,omitempty"`
	CreatedAt   time.Time         `json:"createdAt" dynamodbav:"createdAt,unixtime"`
}
