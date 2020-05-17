package hlfq

import (
	"encoding/json"
	"time"

	cryptorand "crypto/rand"

	"github.com/oklog/ulid"
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

	r.Init(invokeInit)

	r.
		Invoke("Push", queuePush, pdef.Struct(newItemSpecParamName, &QueueItemSpec{})). // 1 struct argument, insert an item to the end of queue (chaincode method name `hlfqueuePush`)
		Invoke("Pop", queuePop).                                                        // 1 struct argument, get the oldes item and delete it from queue
		Invoke("ListItems", queueListItems).
		Invoke("AttachData", queueAttachData, pdef.String(keyParamName), pdef.Bytes(attachedDataParamName)).
		Query("Select", queueSelect, pdef.String(selectMethodParam))

	return router.NewChaincode(r)
}

// **
// ** Chaincode methods **
// **

const queueItemKeyPrefix = "QueueItem"

const (
	newItemSpecParamName  = "newItemSpec"
	popMethodParam        = "extractedItem"
	selectMethodParam     = "selectQueue"
	keyParamName          = "itemKey"
	attachedDataParamName = "attachedData"
)

func invokeInit(c router.Context) (interface{}, error) {
	// no state init required
	return nil, nil
}

// queuePush adds an item after last queue item
func queuePush(c router.Context) (interface{}, error) {
	entropy := cryptorand.Reader
	id, err1 := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err1 != nil {
		return nil, errors.Wrap(err1, "failed generate item UID")
	}
	// getTxTimestamp() - time when transaction proposial was created

	spec := c.Param(newItemSpecParamName).(QueueItemSpec)
	// creare queueItem
	item := &QueueItem{
		ID:        id.String(),
		From:      spec.From,
		To:        spec.To,
		Amount:    spec.Amount,
		ExtraData: spec.ExtraData,
		CreatedAt: time.Now().UTC(),
	}

	return item, c.State().Insert(item)
}

// queuePop read and delete the first queue item (the oldest, FIFO)
func queuePop(c router.Context) (extractedItem interface{}, err error) {

	//return item, c.State().Delete(item)
	return extractedItem, err
}

// attcahes extra data to the item specified by key, returns error if key not exists
// arg1 -> attachDataMethodParamKey string
// arg2 -> attachDataMethodParamData []bytes
func queueAttachData(c router.Context) (interface{}, error) {
	//p := c.Param(attachDataMethodParam)
	//item := c.State().Get()
	return nil, nil
}

// returns error if key not exists
func queueMoveAfter(akey string, afterKey string) (interface{}, error) {
	return nil, nil
}

// returns error if key not exists
func queueMoveBefore(akey string, beforeKey string) (interface{}, error) {
	return nil, nil
}

// Select get elemets specified by CouchDB query
// arg1 =`queryString` - query in CouchDB syntax
// returns error query syntax is invalid
func queueSelect(c router.Context) (interface{}, error) {
	return nil, nil
}

// read and return all queue items as list
func queueListItems(c router.Context) (interface{}, error) {
	// TODO: возврващать всегда в одном порядке сортировать по ID (ULID)
	return c.State().List(queueItemKeyPrefix, &QueueItem{})
}

func queueListItemsSorted(c router.Context) (interface{}, error) {
	// сортировать по ID (т.к. это ULID отсортируются как по времени, первый будет самый старый)
	queryString := `{
        "selector": {},
        "sort": [
            {"id": "asc"}
        ]
	}`

	iter, err1 := c.Stub().GetQueryResult(queryString)
	if err1 != nil {
		return nil, errors.Wrap(err1, "failed to GetQueryResult") // TODO: тут падает в mock-тесте с `not implemended`
	}
	defer iter.Close()

	items := []interface{}{}
	for iter.HasNext() {
		kvResult, err := iter.Next()
		if err != nil {
			return nil, errors.Wrap(err1, "fetch error")
		}
		item := QueueItem{}
		err2 := json.Unmarshal(kvResult.Value, &item)
		if err2 != nil {
			return nil, errors.Wrap(err1, "value unmarshal error")
		}
		items = append(items, item)
	}

	return items, nil
}
