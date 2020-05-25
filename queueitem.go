package hlfq

import (
	"fmt"
	"reflect"
	"time"

	"github.com/oklog/ulid/v2"
)

const queuePointerTypeName = "queuePointer"

// QueuePointer holds a key pointing to another state
type QueuePointer struct {
	PointerName string
	PointerKey  []string
}

// NewQueuePointer creates new QueuePointer (for Head or Tail)
func NewQueuePointer(name string) *QueuePointer {
	return &QueuePointer{PointerName: name}
}

// NewQueueHeadPointer creates a QueuePointer for HEAD
func NewQueueHeadPointer() *QueuePointer {
	return &QueuePointer{PointerName: "HeadPointer"}
}

// NewQueueTailPointer creates a QueuePointer for TAIL
func NewQueueTailPointer() *QueuePointer {
	return &QueuePointer{PointerName: "TailPointer"}
}

// Key for QueuePointer entry in chaincode state
func (qp QueuePointer) Key() ([]string, error) {
	s := []string{queuePointerTypeName}
	s = append(s, qp.PointerName)
	return s, nil
}

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
	PrevKey   []string  `json:"prevKey"`
	NextKey   []string  `json:"nextKey"`
	ID        ulid.ULID `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    int       `json:"amount"`
	ExtraData []byte    `json:"extra"`

	UpdatedTime time.Time `json:"updated_time"` // set by chaincode method
}

// Key for QueueItem entry in chaincode state
func (qi QueueItem) Key() ([]string, error) {
	return []string{queueItemKeyPrefix, qi.ID.String()}, nil
}

func (qi QueueItem) hasNext() bool {
	fmt.Println("== hasNext ==")
	fmt.Printf("qi=%+v\n", qi)
	return !reflect.DeepEqual(qi.NextKey, EmptyItemPointerKey)
}

func (qi QueueItem) hasPrev() bool {
	fmt.Println("== hasPrev ==")
	fmt.Printf("qi=%+v\n", qi)
	return !reflect.DeepEqual(qi.PrevKey, EmptyItemPointerKey)
}
