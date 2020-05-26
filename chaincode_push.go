package hlfq

import (
	"time"

	cryptorand "crypto/rand"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

// **
// ** Chaincode method **
// **

const newItemSpecParam = "newItemSpec"

// queuePush adds an item after last queue item
func queuePush(c router.Context) (interface{}, error) {
	spec := c.Param(newItemSpecParam).(QueueItemSpec)
	// getTxTimestamp() - time when transaction proposial was created
	t, _ := c.Time()                     // tx time // TODO: handle get txt time error
	curItem, _ := makeQueueItem(spec, t) // TODO: handle assign errors
	curItemKey, _ := curItem.Key()       // TODO: handle read error

	tailPresent, _ := hasTail(c) // TODO: handle read error
	if tailPresent {
		tailItem, _ := getTailItem(c) // TODO: handle read error
		tailItem.NextKey = curItemKey // TAIL->CUR
		tailKey, _ := tailItem.Key()  // TODO: handle read error
		curItem.PrevKey = tailKey     // TAIL<-CUR
		// update prvious tail item
		c.State().Put(tailItem) // TODO: handle errors
	}

	// set CUR as tail / replce tail kay with new one
	storeTailKey(c, curItemKey) // TAIL = CUR

	// UPDATE Head key if head not set
	headPresent, _ := hasHead(c) // TODO: handle error
	if !headPresent {
		// c.Logger().Debug("*** headNotPresent")
		// set head pointer to CUR
		storeHeadKey(c, curItemKey) // TODO: handle store write error
	}
	// printout updated states
	// h1, _ := readHeadItemKey(c)
	// t1, _ := readTailItemKey(c)
	// fmt.Printf("queuePush head: %+v\n", h1)
	// fmt.Printf("queuePush tail: %+v\n", t1)

	// insert return an error if item already exists
	return curItem, c.State().Insert(curItem)
}

func makeQueueItem(spec QueueItemSpec, t time.Time) (*QueueItem, error) {
	entropy := cryptorand.Reader
	id, err1 := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err1 != nil {
		return nil, errors.Wrap(err1, "failed generate item UID")
	}
	// data for chaincode state
	item := &QueueItem{
		ID:          id,
		From:        spec.From,
		To:          spec.To,
		Amount:      spec.Amount,
		ExtraData:   spec.ExtraData,
		UpdatedTime: t,
		NextKey:     EmptyItemPointerKey,
		PrevKey:     EmptyItemPointerKey,
	}
	return item, nil

}
