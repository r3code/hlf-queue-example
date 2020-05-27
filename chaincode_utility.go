package hlfq

import (
	"fmt"
	"reflect"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

//
// *** Utilty methods ***
//

// IsEmpty shows if queue not empty
/*func IsEmpty(c router.Context) (empty bool, err error) {
	headKey, err1 := readHeadItemKey(c)
	if err1 != nil {
		return empty, errors.Wrap(err1, "check queue empty")
	}
	if headKey != EmptyItemPointerKey {
		return false, nil
	}
	tailKey, err2 := readTailItemKey(c)
	if err2 != nil {
		return empty, errors.Wrap(err2, "check queue empty")
	}
	if tailKey != EmptyItemPointerKey {
		return false, nil
	}
	return true, nil
}*/

// read current key of head item
func readHeadItemKey(c router.Context) (headKey []string, err error) {
	pointerKey, _ := NewQueueHeadPointer().Key()
	headPointer, err := readQueuePointer(c, pointerKey)
	if err != nil {
		return headKey, errors.Wrap(err, "failed to read key of a head item")
	}
	//fmt.Printf("*HEAD* %+v ", headPointer)
	headKey = headPointer.PointerKey
	return headKey, nil
}

// read current key of tail item
func readTailItemKey(c router.Context) (tailKey []string, err error) {
	pointerKey, _ := NewQueueTailPointer().Key()
	tailPointer, err := readQueuePointer(c, pointerKey)
	if err != nil {
		return tailKey, errors.Wrap(err, "failed to read key of a tail item")
	}
	//fmt.Printf("****TAIL****** %+v ", tailPointer)

	return tailPointer.PointerKey, nil
}

// replace a tail pointer with itemKey
func setHeadPointerTo(c router.Context, itemKey []string) (err error) {
	fmt.Printf("\n::--STORE HEAD: %v\n\n", itemKey)
	headPointer := NewQueueHeadPointer()
	headPointer.PointerKey = itemKey
	if err := c.State().Put(headPointer); err != nil {
		return errors.Wrap(err, "failed to update queue head pointer")
	}
	return nil
}

// replace a tail pointer with itemKey
func setTailPointerTo(c router.Context, itemKey []string) (err error) {
	fmt.Printf("\n--::STORE TAIL: %v\n\n", itemKey)
	tailPointer := NewQueueTailPointer()
	tailPointer.PointerKey = itemKey

	if err := c.State().Put(tailPointer); err != nil {
		return errors.Wrap(err, "failed to update queue head pointer")
	}
	return nil
}

func getHeadItem(c router.Context) (headItem QueueItem, err error) {
	headItemKey, err := readHeadItemKey(c)
	if err != nil {
		return headItem, err
	}
	headItem, err = readQueueItem(c, headItemKey)
	if err != nil {
		return headItem, errors.Wrap(err, "failed to load head item")
	}
	return headItem, nil
}

func getTailItem(c router.Context) (tailItem QueueItem, err error) {
	tailItemKey, err := readTailItemKey(c)
	if err != nil {
		return tailItem, err
	}

	if isKeyEmpty(tailItemKey) {
		return tailItem, errors.New("Empty queue")
	}

	tailItem, err = readQueueItem(c, tailItemKey)
	if err != nil {
		return tailItem, errors.Wrap(err, "failed to load tail item")
	}
	return tailItem, nil
}

func hasTail(c router.Context) (bool, error) {
	tailKey, err := readTailItemKey(c)
	if err != nil {
		return false, err
	}
	if isKeyEmpty(tailKey) {
		return false, nil
	}
	return true, nil
}

func hasHead(c router.Context) (bool, error) {
	headKey, err := readHeadItemKey(c)
	if err != nil {
		return false, err
	}
	if isKeyEmpty(headKey) {
		return false, nil
	}
	return true, nil
}

func isKeyEmpty(key []string) bool {
	return reflect.DeepEqual(key, EmptyItemPointerKey)
}

func readQueueItem(c router.Context, itemKey []string) (item QueueItem, err error) {
	res, err := c.State().Get(itemKey, &QueueItem{})
	if err != nil {
		return item, errors.Wrapf(err, "failed to read QueueItem with key '%v'", itemKey)
	}
	item = res.(QueueItem)
	return item, nil
}

func readQueueItemByID(c router.Context, itemIDStr string) (item QueueItem, err error) {
	id, err := ulid.ParseStrict(itemIDStr)
	if err != nil {
		return item, errors.Wrap(err, "invalid ULID string passed")
	}
	itemForKey := QueueItem{ID: id}
	itemKey, _ := itemForKey.Key()
	res, err := c.State().Get(itemKey, &QueueItem{})
	if err != nil {
		return item, errors.Wrapf(err, "failed to read QueueItem with ID '%s'", itemIDStr)
	}
	item = res.(QueueItem)
	return item, nil
}

func readQueuePointer(c router.Context, key []string) (pointerItem QueuePointer, err error) {
	res, err := c.State().Get(key, &QueuePointer{})
	if err != nil {
		return pointerItem, errors.Wrapf(err, "failed to read QueuePointer with key '%v'", key)
	}
	pointerItem = res.(QueuePointer)
	return pointerItem, nil
}

// connect left item to right
func connectItems(c router.Context, leftIDStr string, rightIDStr string) (err error) {
	leftItem, err := readQueueItemByID(c, leftIDStr)
	if err != nil {
		return errors.Wrapf(err, "failed load leftItem ID '%s'", leftIDStr)
	}
	leftItemKey, _ := leftItem.Key()
	rightItem, err := readQueueItemByID(c, rightIDStr)
	if err != nil {
		return errors.Wrapf(err, "failed load rightItem ID '%s'", rightIDStr)
	}
	rightItemKey, _ := rightItem.Key()
	leftItem.NextKey = rightItemKey
	rightItem.PrevKey = leftItemKey
	err = c.State().Put(leftItem)
	if err != nil {
		return errors.Wrap(err, "faild to save leftItem")
	}
	c.State().Put(rightItem)
	if err != nil {
		return errors.Wrap(err, "faild to save rightItem")
	}
	return nil
}

// Unlinks the item from neighbours, and returns item data (A<->B<->C => [B], A<->C).
// It does not delete the state from the ledger, only disconnects the item from the Prev and Next item,
// PrevItem.NextKey replaced to item.NextKey,
// NextItem.PrevKey replaced to item.PrevKey
// Updates Head and Tail pointer if item is a head or tail item
func cutItem(c router.Context, itemIDStr string) (item QueueItem, err error) {
	item, err = readQueueItemByID(c, itemIDStr)
	if err != nil {
		return item, errors.Wrapf(err, "failed load item ID '%s'", itemIDStr)
	}

	// check if item is a Head, so we need to replace HeadPointer
	if isHeadPointsTo(c, item) {
		// move head pointer to next item (list=X[head]<->Y => list=Y[Head], cut=X)
		setHeadPointerTo(c, item.NextKey) // TODO: handle error
		headItem2, _ := getHeadItem(c)    // TODO: handle error
		fmt.Printf("***NEW HeadID=%s\n", headItem2.ID.String())
	}

	// check if item is a Tail, so we need to replace TailPointer

	if isTailPointsTo(c, item) {
		// set tail pointer to prevous item (list=X->Y[Tail] => list=X[Tail], cut=Y)
		setTailPointerTo(c, item.PrevKey) // TODO: handle error
		tailItem2, _ := getTailItem(c)
		fmt.Printf("--> NEW TailID=%s\n", tailItem2.ID.String())
	}

	// prev <- item -> next
	var prevItem, nextItem QueueItem
	// update PrevID of an item after targetItem if present
	fmt.Printf("\n   *** item_TO_CUT ***\n%+v\n\n", item)
	if item.hasPrev() {
		fmt.Println("***HAS PREV")
		prevItem, err = readQueueItem(c, item.PrevKey)
		if err != nil {
			return item, errors.Wrapf(err, "failed load prev item for ID '%s'", itemIDStr)
		}
		fmt.Printf("*** prevItem.OLD=%+v\n", prevItem)
		prevItem.NextKey = item.NextKey
		// save updated prevItem
		c.State().Put(prevItem) // TODO: handle error
		fmt.Printf("--> prevItem.NEW=%+v\n", prevItem)
	}
	// update NextID of an item before targetItem if present
	if item.hasNext() {
		// fmt.Println("***HAS NEXT")
		nextItem, err = readQueueItem(c, item.NextKey)
		if err != nil {
			return item, errors.Wrapf(err, "failed load next item for ID '%s'", itemIDStr)
		}
		fmt.Printf("*** nextItem.OLD=%+v\n", nextItem)
		nextItem.PrevKey = item.PrevKey
		// save updated nextItem
		c.State().Put(nextItem) // TODO: handle error

		fmt.Printf("*** nextItem.NEW=%+v\n", nextItem)
	}
	return item, nil
}

func isHeadPointsTo(c router.Context, item QueueItem) bool {
	headItem, _ := getHeadItem(c) // TODO: handle error
	return headItem.ID.Compare(item.ID) == 0
}

func isTailPointsTo(c router.Context, item QueueItem) bool {
	tailItem, _ := getTailItem(c) // TODO: handle error
	return tailItem.ID.Compare(item.ID) == 0
}
