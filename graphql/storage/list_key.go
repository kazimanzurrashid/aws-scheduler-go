package storage

import "time"

type ListKey struct {
	ID     string    `json:"id" dynamodbav:"id"`
	DueAt  time.Time `json:"dueAt" dynamodbav:"dueAt,unixtime"`
	Status string    `json:"status,omitempty" dynamodbav:"status,omitempty"`
}
