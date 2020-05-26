package hlfq

import (
	"fmt"

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

	// cut item and reconnect neighbours
	item, err := cutItem(c, itemIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to cut item ID '%s'", itemIDStr)
	}
	itemKey, _ := item.Key()
	// reset links
	item.PrevKey = EmptyItemPointerKey
	item.NextKey = EmptyItemPointerKey

	fmt.Printf("*** CUT_ITEM=%+v", item)

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
	// connect: afterItem -[next]-> item
	afterItem.NextKey = itemKey

	// connect: afterItem <-[prev]- item
	item.PrevKey = afterItemKey

	// save link update of item. Item now between after and afterNext items
	c.State().Put(item)      // TODO: handle error
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
