package hlfq

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

// **
// ** Chaincode methods **
// **

// EmptyItemPointerKey HEAD or TAIL value when queue is empty
const EmptyItemPointerKey = "*EMPTY*"

func invokeInitLedger(c router.Context) (interface{}, error) {
	var headPointer *QueuePointer = NewQueueHeadPointer()
	headPointer.PointerKey = []string{EmptyItemPointerKey}
	var tailPointer *QueuePointer = NewQueueTailPointer()
	tailPointer.PointerKey = []string{EmptyItemPointerKey}

	// init place to store a head pointer
	if err := c.State().Insert(headPointer); err != nil {
		return nil, errors.Wrap(err, "failed to init head pointer store")
	}
	// init place to store a tail pointer
	if err := c.State().Insert(tailPointer); err != nil {
		return nil, errors.Wrap(err, "failed to init tail pointer store")
	}
	return nil, nil
}
