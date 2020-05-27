package hlfq

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

// queueMoveAfter cuts item and puts it after specified item ID.
// returns an updated item (should have updated prev/next links)
//   or error if itemID or afterItemID not exists
// arg1 -> itemID string (ULID String)
// arg2 -> afterItemID string (ULID String)
func queueMoveAfter(c router.Context) (interface{}, error) {
	itemIDStr := c.ParamString(itemIDParam)
	afterItemIDStr := c.ParamString(afterItemIDParam)
	if itemIDStr == afterItemIDStr {
		return nil, errors.New("Can not move an item after itself")
	}

	// cut item and reconnect neighbours
	item, err := cutItem(c, itemIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to cut item ID '%s'", itemIDStr)
	}
	itemKey, _ := item.Key()

	// fmt.Printf("\n\n** queueMoveAfter :: CUT_ITEM=%s\n\n", item.String())
	// reset links
	item.PrevKey = EmptyItemPointerKey
	item.NextKey = EmptyItemPointerKey

	afterItem, err := readQueueItemByID(c, afterItemIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed load afterItem ID '%s'", afterItemIDStr)
	}

	afterItemKey, _ := afterItem.Key()
	if afterItem.hasNext() {
		// need to update Prev in afterNext item
		// after <-> afterNext
		afterItemNext, err := readQueueItem(c, afterItem.NextKey)
		if err != nil {
			return nil, errors.Wrapf(err, "failed load afterItemNext Key='%v'", afterItem.NextKey)
		}
		afterItemNextKey, _ := afterItemNext.Key()
		// connect: item <-[prev]- afterNext
		afterItemNext.PrevKey = itemKey
		// connect: item -[next]-> afterNext
		item.NextKey = afterItemNextKey
		// save link update of afterItemNext
		c.State().Put(afterItemNext) // TODO: handle error
	}

	// Update the Tail pointer if we paste after the tail item
	// check if item is a Tail, so we need to replace TailPointer

	if isTailPointsTo(c, afterItem) { // pasting after tail item
		// item now is new tail
		setTailPointerTo(c, itemKey) // TODO: handle error
		// tailItem2, _ := getTailItem(c)
		// fmt.Printf("queueMoveAfter::--> NEW TailID=%s\n", tailItem2.ID.String())
	}

	// connect: afterItem -[next]-> item
	afterItem.NextKey = itemKey

	// connect: afterItem <-[prev]- item
	item.PrevKey = afterItemKey

	// save link update of item. Item now between after and afterNext items
	c.State().Put(item)      // TODO: handle error
	c.State().Put(afterItem) // TODO: handle error

	return item, nil
}

// queueMoveBefore cuts item and puts it before specified item ID.
// returns an updated item (should have updated prev/next links)
//   error if itemID or afterItemID not exists
// arg1 -> itemID string (ULID String)
// arg2 -> beforeItemID string (ULID String)
func queueMoveBefore(c router.Context) (interface{}, error) {
	itemIDStr := c.ParamString(itemIDParam)
	beforeItemIDStr := c.ParamString(beforeItemIDParam)
	if itemIDStr == beforeItemIDStr {
		return nil, errors.New("Can not move an item before itself")
	}

	// cut item and reconnect neighbours
	item, err := cutItem(c, itemIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to cut item ID '%s'", itemIDStr)
	}
	itemKey, _ := item.Key()

	// fmt.Printf("\n\n** queueMoveBefore:: CUT_ITEM=%s\n\n", item.String())
	// reset links
	item.PrevKey = EmptyItemPointerKey
	item.NextKey = EmptyItemPointerKey

	beforeItem, err := readQueueItemByID(c, beforeItemIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed load beforeItem ID '%s'", beforeItemIDStr)
	}

	beforeItemKey, _ := beforeItem.Key()
	if beforeItem.hasPrev() {
		// need to update Next in beforePrev item
		// after <-> afterNext
		beforeItemPrev, err := readQueueItem(c, beforeItem.PrevKey)
		if err != nil {
			return nil, errors.Wrapf(err, "failed load beforeItemNext Key='%v'", beforeItem.PrevKey)
		}
		beforeItemPrevKey, _ := beforeItemPrev.Key()
		// connect: item -[next]-> beforePrev
		// make a chain: beforeItemPrev -> item -> beforeItem
		beforeItemPrev.NextKey = itemKey
		// connect: item <-[prev]- beforePrev
		item.PrevKey = beforeItemPrevKey
		// save link update of beforeItemPrev
		c.State().Put(beforeItemPrev) // TODO: handle error
	}

	// Update the Head pointer if we paste before the head item
	// check if item is a Head, so we need to replace HeadPointer
	if isHeadPointsTo(c, beforeItem) { // TODO: handle error
		// fmt.Println("queueMoveBefore:: beforeItem is HEAD --->[...]")
		// fmt.Printf("queueMoveBefore::***OLD HeadID=%s\n", beforeItem.ID.String())
		// set head pointer to next item (list=X[head]<->Y => list=Y[Head], cut=X)
		setHeadPointerTo(c, itemKey) // TODO: handle error
		// headItem2, _ := getHeadItem(c) // TODO: handle error
		// fmt.Printf("queueMoveBefore::***NEW HeadID=%s\n", headItem2.ID.String())
	}

	// connect:  <-[prev]- item
	beforeItem.PrevKey = itemKey

	// connect: item -[next]-> beforeItem
	item.NextKey = beforeItemKey

	// save link update of item. Item now between beforeItemPrev and beforeItem items
	c.State().Put(item)       // TODO: handle error
	c.State().Put(beforeItem) // TODO: handle error

	return item, nil
}
