package storage

type ListInput struct {
	Status   string     `json:"status,omitempty"`
	DueAt    *DateRange `json:"dueAt,omitempty"`
	StartKey *ListKey   `json:"startKey,omitempty"`
	Limit    int64      `json:"limit"`
}
