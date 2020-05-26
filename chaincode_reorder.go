package hlfq

import (
	"github.com/s7techlab/cckit/router"
)

// returns error if itemID or afterItemID not exists
// arg1 -> itemID string (ULID String)
// arg2 -> afterItemID string (ULID String)
func queueMoveAfter(c router.Context) (interface{}, error) {
	// TODO: переставить элемент с указанным ключем после заданного
	// 1. Получить ключи
	// 2. Получить элементы item, item.Prev, item.Next
	// 3. Получить элементы after, after.Next
	// GetItem(after.Next).Prev = item
	// 4. afterNext := after.Next; after.Next -> item
	// 5. item.Prev -> after
	// 6. item.Next = afterNext
	// 7.
	return nil, nil
}

//  returns error if itemID or afterItemID not exists
// arg1 -> itemID string (ULID String)
// arg2 -> beforeKey string (ULID String)
func queueMoveBefore(c router.Context) (interface{}, error) {

	return nil, nil
}
