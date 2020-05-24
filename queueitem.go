package hlfq

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// QueueItemSpec chaincode method argument
type QueueItemSpec struct {
	ID        ulid.ULID
	From      string
	To        string
	Amount    int
	ExtraData []byte
}

// QueueItem struct for chaincode state
type QueueItem struct {
	ID        ulid.ULID `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    int       `json:"amount"`
	ExtraData []byte    `json:"extra"`

	UpdatedTime time.Time `json:"updated_time"` // set by chaincode method
}

// Key for QueueItem entry in chaincode state
func (c QueueItem) Key() ([]string, error) {
	return []string{queueKeyPrefix, c.ID.String()}, nil
}
