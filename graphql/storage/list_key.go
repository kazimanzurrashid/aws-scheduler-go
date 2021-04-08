package storage

type ListKey struct {
	ID     string  `json:"id" dynamodbav:"id"`
	DueAt  *int64  `json:"dueAt" dynamodbav:"dueAt"`
	Status *string `json:"status,omitempty" dynamodbav:"status,omitempty"`
}
