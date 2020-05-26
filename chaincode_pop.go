package hlfq

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

// queuePop read and delete the first queue item (the oldest, FIFO)
func queuePop(c router.Context) (extractedItem interface{}, err error) {
	headPresent, _ := hasHead(c) // TODO: handle error
	if !headPresent {
		return extractedItem, errors.New("Empty queue")
	}

	headKey, _ := readHeadItemKey(c)                   // TODO: handle error
	resHead, _ := c.State().Get(headKey, &QueueItem{}) // TODO: handle error
	headItem := resHead.(QueueItem)

	nextKey := EmptyItemPointerKey
	// remove Prev link from nextItem if it exists
	if headItem.hasNext() { // it's not a tail
		// получить следующий элемент
		nextKey = headItem.NextKey
		// получить из State nextItem
		resNext, _ := c.State().Get(nextKey, &QueueItem{}) // TODO: handle error
		nextItem := resNext.(QueueItem)
		// remove PrevKey from nextItem
		nextItem.PrevKey = EmptyItemPointerKey
		// save updated nextItem
		c.State().Put(nextItem) // TODO: handle error
	}
	storeHeadKey(c, nextKey)
	if isKeyEmpty(nextKey) { // reached a TailItem
		storeTailKey(c, nextKey)
	}
	// remove extracted item from state
	c.State().Delete(headKey) // TODO: handle error

	//return item, c.State().Delete(item)
	extractedItem = headItem
	return extractedItem, nil
}
