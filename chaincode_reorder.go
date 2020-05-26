package hlfq

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

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

	// cut item and connect neighbours
	// A<->B<->C => [B], A<->C
	item, err := cutItem(c, itemIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to cut item ID '%s'", itemIDStr)
	}

	afterItem, err := readQueueItemByID(c, afterItemIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed load afterItem ID '%s'", afterItemIDStr)
	}
	// after <-> afterNext
	afterItemNext, err := readQueueItem(c, afterItem.NextKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed load afterItemNext Key='%v'", afterItem.NextKey)
	}
	// put item next to afterItem
	itemKey, _ := item.Key()
	// connect: item <-[prev]- afterNext
	afterItemNext.PrevKey = itemKey
	// save link update of afterItemNext
	c.State().Put(afterItemNext) // TODO: handle error
	afterItemNextKey, _ := afterItemNext.Key()
	// connect: item -[next]-> afterNext
	item.NextKey = afterItemNextKey
	afterItemKey, _ := afterItem.Key()
	item.PrevKey = afterItemKey
	// save link update of item. Item now between after and afterNext items
	c.State().Put(item) // TODO: handle error
	// connect: afterItem -[next]-> item
	afterItem.NextKey = itemKey
	c.State().Put(afterItem) // TODO: handle error
	return item, nil
}

// returns an updated item (should have updated prev/next links)
//   error if itemID or afterItemID not exists
// arg1 -> itemID string (ULID String)
// arg2 -> beforeItemID string (ULID String)
func queueMoveBefore(c router.Context) (interface{}, error) {
	// itemIDStr := c.ParamString(itemIDParam)
	// beforeItemIDStr := c.ParamString(beforeItemIDParam)
	// if itemIDStr == beforeItemIDStr {
	// 	return nil, errors.New("Can not move an item before itself")
	// }

	// item, err := cutItem(c, itemIDStr)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "failed load item ID '%s'", itemIDStr)
	// }

	// beforeItem, err := readQueueItemByID(c, beforeItemIDStr)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "failed load beforeItem ID '%s'", beforeItemIDStr)
	// }

	return nil, nil
}

// Unlinks the item from neighbours, and returns item data (A<->B<->C => [B], A<->C).
// It does not delete the state from the ledger, only disconnects the item from the Prev and Next item,
// PrevItem.NextKey replaced to item.NextKey,
// NextItem.PrevKey replaced to item.PrevKey
func cutItem(c router.Context, itemIDStr string) (cutItem QueueItem, err error) {
	cutItem, err = readQueueItemByID(c, itemIDStr)
	if err != nil {
		return cutItem, errors.Wrapf(err, "failed load item ID '%s'", itemIDStr)
	}
	// prev <- item -> next
	prevKey := cutItem.PrevKey
	nextKey := cutItem.NextKey
	// remove links from the item being cut
	cutItem.NextKey = EmptyItemPointerKey
	cutItem.PrevKey = EmptyItemPointerKey
	// update PrevID of an item after targetItem if present
	if cutItem.hasPrev() {
		prevItem, err1 := readQueueItem(c, cutItem.NextKey)
		if err1 != nil {
			return cutItem, errors.Wrapf(err1, "failed load prev item for ID '%s'", itemIDStr)
		}
		prevItem.NextKey = nextKey
		// save updated prevItem
		c.State().Put(prevItem) // TODO: handle error
	}
	// update NextID of an item before targetItem if present
	if cutItem.hasNext() {
		nextItem, err2 := readQueueItem(c, cutItem.NextKey)
		if err2 != nil {
			return cutItem, errors.Wrapf(err2, "failed load next item for ID '%s'", itemIDStr)
		}
		nextItem.PrevKey = prevKey
		// save updated nextItem
		c.State().Put(nextItem) // TODO: handle error
	}
	return cutItem, nil
}
