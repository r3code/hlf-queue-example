package hlfq

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
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
	queueItemKeyPrefix     = "queueItemKey"
	newItemSpecParam       = "newItemSpec"
	extractedItemParam     = "extractedItem"
	selectQueryStringParam = "queryString"
	itemIDParam            = "itemID"
	initQueueNameParam     = "queueName"
	attachedDataParam      = "attachedData"
)

// Select get elemets specified by CouchDB query
// arg1 =`queryString` - query in `expr` syntax
// returns error query syntax is invalid
func queueSelect(c router.Context) (interface{}, error) {
	queryStr := c.ParamString(selectQueryStringParam)
	res, err := queueListItems(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read queue for Select")
	}
	items := res.([]QueueItem)
	type env struct {
		QueueItems []QueueItem
	}
	queryStr = fmt.Sprintf("filter(QueueItems, %s)", queryStr)
	program, err := expr.Compile(queryStr, expr.Env(env{}))
	if err != nil {
		return nil, errors.Wrap(err, "queryString parse error")
	}
	progEnv := env{
		QueueItems: items,
	}

	filteredItems, err := expr.Run(program, progEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed filter operation")
	}
	// fmt.Printf("***********filtered=%+v\n", filteredItems)
	return filteredItems, nil
}

func queueListItems(c router.Context) (interface{}, error) {
	// we can raplace realization to any of queueListItems*
	return queueListItemsItarated(c)
}
