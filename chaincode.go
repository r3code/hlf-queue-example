package hlfq

import (
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	pdef "github.com/s7techlab/cckit/router/param"
)

const (
	queueItemKeyPrefix = "queueItemKey"
	itemIDParam        = "itemID"
	afterItemIDParam   = "afterItemID"
	beforeItemIDParam  = "beforeItemID"
)

// New inits a chaincode, adds chaincode methods to the rourer
// All methods allow access to anyone
func New() *router.Chaincode {
	r := router.New("hlfq") // also initialized logger with "hlfq_*" prefix

	// Method for debug chaincode state
	debug.AddHandlers(r, "debug", owner.Only)

	r.Init(invokeInitLedger) // no params

	r.
		Invoke("Push", queuePush, pdef.Struct(newItemSpecParam, &QueueItemSpec{})).
		Invoke("Pop", queuePop).
		Invoke("ListItems", queueListItems).
		Invoke("AttachData", queueAttachData, pdef.String(itemIDParam), pdef.Bytes(attachedDataParam)).
		Invoke("MoveAfter", queueMoveAfter, pdef.String(itemIDParam), pdef.String(afterItemIDParam)).
		Invoke("MoveBefore", queueMoveBefore, pdef.String(itemIDParam), pdef.String(beforeItemIDParam)).
		Query("Select", queueSelect, pdef.String(selectQueryStringParam))

	return router.NewChaincode(r)
}
