package storage

type List struct {
	Schedules []*Schedule `json:"schedules"`
	NextKey   *ListKey    `json:"nextKey,omitempty"`
}
