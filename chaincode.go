package hlfq

import (
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	pdef "github.com/s7techlab/cckit/router/param"
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
		Query("Select", queueSelect, pdef.String(selectQueryStringParam))

	return router.NewChaincode(r)
}

// **
// ** Chaincode methods **
// **

const (
	queueItemKeyPrefix = "queueItemKey"
	itemIDParam        = "itemID"
)

func queueListItems(c router.Context) (interface{}, error) {
	// we can raplace realization to any of queueListItems*
	return queueListItemsItarated(c)
}
