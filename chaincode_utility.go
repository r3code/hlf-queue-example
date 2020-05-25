package hlfq

import (
	"reflect"

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
	res, err := c.State().Get(NewQueueHeadPointer(), &QueuePointer{})
	if err != nil {
		return headKey, errors.Wrap(err, "failed to read key of a head item")
	}
	headPointer := res.(QueuePointer)
	//fmt.Printf("*HEAD* %+v ", headPointer)
	headKey = headPointer.PointerKey
	return headKey, nil
}

// read current key of tail item
func readTailItemKey(c router.Context) (tailKey []string, err error) {
	res, err := c.State().Get(NewQueueTailPointer(), &QueuePointer{})
	if err != nil {
		return tailKey, errors.Wrap(err, "failed to read key of a tail item")
	}
	tailPointer := res.(QueuePointer)
	//fmt.Printf("****TAIL****** %+v ", tailPointer)

	return tailPointer.PointerKey, nil
}

// replace a tail pointer with itemKey
func storeHeadKey(c router.Context, itemKey []string) (err error) {
	headPointer := NewQueueHeadPointer()
	headPointer.PointerKey = itemKey
	if err := c.State().Put(headPointer); err != nil {
		return errors.Wrap(err, "failed to update queue head pointer")
	}
	return nil
}

// replace a tail pointer with itemKey
func storeTailKey(c router.Context, itemKey []string) (err error) {
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

	res, err := c.State().Get(headItemKey, &QueueItem{})
	if err != nil {
		return headItem, errors.Wrap(err, "failed to load head item")
	}
	headItem = res.(QueueItem)
	return headItem, nil
}

func getTailItem(c router.Context) (tailItem QueueItem, err error) {
	tailItemKey, err := readTailItemKey(c)
	if err != nil {
		return tailItem, err
	}

	if reflect.DeepEqual(tailItemKey, EmptyItemPointerKey) {
		return tailItem, errors.New("Empty queue")
	}

	res, err := c.State().Get(tailItemKey, &QueueItem{})
	if err != nil {
		return tailItem, errors.Wrap(err, "failed to load tail item")
	}
	tailItem = res.(QueueItem)
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
