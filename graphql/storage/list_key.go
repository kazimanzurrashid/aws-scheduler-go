package storage

import "time"

type ListKey struct {
	ID     string    `json:"id" dynamodbav:"id"`
	DueAt  time.Time `json:"dueAt,omitempty" dynamodbav:"dueAt,unixtime,omitempty"`
	Status string    `json:"status,omitempty" dynamodbav:"status,omitempty"`
}
