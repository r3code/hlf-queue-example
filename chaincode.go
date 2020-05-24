package hlfq

import (
	"encoding/json"
	"sort"
	"time"

	cryptorand "crypto/rand"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	pdef "github.com/s7techlab/cckit/router/param"
)

// EmptyItemPointer a valus to store in head and tail pointer when queue is empty
const EmptyItemPointer = "*NIL*"

// New inits a chaincode, adds chaincode methods to the rourer
// All methods allow access to anyone
func New() *router.Chaincode {
	r := router.New("hlfq") // also initialized logger with "hlfq_*" prefix

	// Method for debug chaincode state
	debug.AddHandlers(r, "debug", owner.Only)

	r.Init(invokeInit)

	r.
		Invoke("Push", queuePush, pdef.Struct(newItemSpecParam, &QueueItemSpec{})). // 1 struct argument, insert an item to the end of queue (chaincode method name `hlfqueuePush`)
		Invoke("Pop", queuePop).                                                    // 1 struct argument, get the oldes item and delete it from queue
		Invoke("ListItems", queueListItemsAsIs).
		Invoke("AttachData", queueAttachData, pdef.String(itemKeyParam), pdef.Bytes(attachedDataParam)).
		Query("Select", queueSelect, pdef.String(selectQueryStringParam))

	return router.NewChaincode(r)
}

// **
// ** Chaincode methods **
// **

const (
	queueKeyPrefix         = "Queue"
	newItemSpecParam       = "newItemSpec"
	extractedItemParam     = "extractedItem"
	selectQueryStringParam = "queryString"
	itemKeyParam           = "itemKey"
	initQueueNameParam     = "queueName"
	attachedDataParam      = "attachedData"
	// a state key name for a tag block hoding a key name of the first queue item
	headStoreKey = queueKeyPrefix + "~HEAD"
	// a state key name for a tag block hoding a key name of the last queue item
	tailStoreKey = queueKeyPrefix + "~TAIL"
)

func invokeInit(c router.Context) (interface{}, error) {
	// init place to store a head pointer
	if err := c.State().Insert(headStoreKey, EmptyItemPointer); err != nil {
		return nil, errors.Wrap(err, "failed to init head pointer store")
	}
	// init place to store a tail pointer
	if err := c.State().Insert(tailStoreKey, EmptyItemPointer); err != nil {
		return nil, errors.Wrap(err, "failed to init tail pointer store")
	}
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
	t, _ := c.Time() // tx time
	spec := c.Param(newItemSpecParam).(QueueItemSpec)
	// creare queueItem
	item := &QueueItem{
		ID:          id,
		From:        spec.From,
		To:          spec.To,
		Amount:      spec.Amount,
		ExtraData:   spec.ExtraData,
		UpdatedTime: t,
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

// queueListItems read and return all queue items as list sorted by ULID stored in ID
func queueListItems(c router.Context) (interface{}, error) {
	res, err := c.State().List(queueKeyPrefix, &QueueItem{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list queue items")
	}
	list := res.([]interface{})

	var items []QueueItem
	for _, item := range list {
		items = append(items, item.(QueueItem))
	}
	// fmt.Printf("queueListItems: unsorted items = %v\n\n", items)
	sort.SliceStable(items, func(i, j int) bool {
		//fmt.Printf("queueListItems: less %v, %v < %v\n\n", (items[i].ID.Compare(items[j].ID) < 0), items[i].ID.String(), items[j].ID.String())
		return items[i].ID.Compare(items[j].ID) < 0
	})
	return items, nil
}

// queueListItemsSorted implemens lising queue by CouchDB rich query
// NOTE: you can not test it by Mock, it doesn't implement GetQueryResult()
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

// queueListItemsAsIs returns queue items in order they retreived from state DB (unexpected)
func queueListItemsAsIs(c router.Context) (interface{}, error) {
	return c.State().List(queueKeyPrefix, &QueueItem{})
}

//
// *** Utilty methods ***
//
func getQueueHeadKey(c router.Context) (tailKey string, err error) {
	res, err := c.State().Get(headStoreKey)
	if err != nil {
		return tailKey, errors.Wrap(err, "failed to read key of a head item")
	}
	tailKey = res.(string)
	return tailKey, nil
}

func getQueueTailKey(c router.Context) (tailKey string, err error) {
	res, err := c.State().Get(tailStoreKey)
	if err != nil {
		return tailKey, errors.Wrap(err, "failed to read key of a tail item")
	}
	tailKey = res.(string)
	return tailKey, nil
}
