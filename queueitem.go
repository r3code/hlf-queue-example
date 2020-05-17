package hlfq

import "time"

// QueueItemSpec chaincode method argument
type QueueItemSpec struct {
	ID        string
	From      string
	To        string
	Amount    int
	ExtraData []byte
}

// QueueItem struct for chaincode state
type QueueItem struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    int    `json:"amount"`
	ExtraData []byte `json:"extra"`

	CreatedAt time.Time `json:"created_at"` // set by chaincode method
}

// Key for QueueItem entry in chaincode state
func (c QueueItem) Key() ([]string, error) {
	return []string{queueItemKeyPrefix, c.ID, c.CreatedAt.String()}, nil
}
