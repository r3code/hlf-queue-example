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

// New inits a chaincode, adds chaincode methods to the rourer
// All methods allow access to anyone
func New() *router.Chaincode {
	r := router.New("hlfq") // also initialized logger with "hlfq_*" prefix

	// Method for debug chaincode state
	debug.AddHandlers(r, "debug", owner.Only)

	r.Init(invokeInitLedger) // no params

	r.
		Invoke("Push", queuePush, pdef.Struct(newItemSpecParam, &QueueItemSpec{})). // 1 struct argument, insert an item to the end of queue (chaincode method name `hlfqueuePush`)
		Invoke("Pop", queuePop).                                                    // 1 struct argument, get the oldes item and delete it from queue
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

// queuePush adds an item after last queue item
func queuePush(c router.Context) (interface{}, error) {
	spec := c.Param(newItemSpecParam).(QueueItemSpec)
	// getTxTimestamp() - time when transaction proposial was created
	t, _ := c.Time()                     // tx time // TODO: handle get txt time error
	curItem, _ := makeQueueItem(spec, t) // TODO: handle assign errors
	curItemKey, _ := curItem.Key()       // TODO: handle read error

	tailPresent, _ := hasTail(c) // TODO: handle read error
	if tailPresent {
		tailItem, _ := getTailItem(c) // TODO: handle read error
		tailItem.NextKey = curItemKey // TAIL->CUR
		tailKey, _ := tailItem.Key()  // TODO: handle read error
		curItem.PrevKey = tailKey     // TAIL<-CUR
		// update prvious tail item
		c.State().Put(tailItem) // TODO: handle errors
	}

	// set CUR as tail / replce tail kay with new one
	storeTailKey(c, curItemKey) // TAIL = CUR

	// UPDATE Head key if head not set
	headPresent, _ := hasHead(c) // TODO: handle error
	if !headPresent {
		// c.Logger().Debug("*** headNotPresent")
		// set head pointer to CUR
		storeHeadKey(c, curItemKey) // TODO: handle store write error
	}
	// printout updated states
	// h1, _ := readHeadItemKey(c)
	// t1, _ := readTailItemKey(c)
	// fmt.Printf("queuePush head: %+v\n", h1)
	// fmt.Printf("queuePush tail: %+v\n", t1)

	// insert return an error if item already exists
	return curItem, c.State().Insert(curItem)
}

func makeQueueItem(spec QueueItemSpec, t time.Time) (*QueueItem, error) {
	entropy := cryptorand.Reader
	id, err1 := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err1 != nil {
		return nil, errors.Wrap(err1, "failed generate item UID")
	}
	// data for chaincode state
	item := &QueueItem{
		ID:          id,
		From:        spec.From,
		To:          spec.To,
		Amount:      spec.Amount,
		ExtraData:   spec.ExtraData,
		UpdatedTime: t,
		NextKey:     EmptyItemPointerKey,
		PrevKey:     EmptyItemPointerKey,
	}
	return item, nil

}

// queuePop read and delete the first queue item (the oldest, FIFO)
func queuePop(c router.Context) (extractedItem interface{}, err error) {
	headPresent, _ := hasHead(c) // TODO: handle error
	if !headPresent {
		return extractedItem, errors.New("Empty queue")
	}

	headKey, _ := readHeadItemKey(c)                   // TODO: handle error
	resHead, _ := c.State().Get(headKey, &QueueItem{}) // TODO: handle error
	headItem := resHead.(QueueItem)

	nextKey := EmptyItemPointerKey
	// remove Prev link from nextItem if it exists
	if headItem.hasNext() { // it's not a tail
		// получить следующий элемент
		nextKey = headItem.NextKey
		// получить из State nextItem
		resNext, _ := c.State().Get(nextKey, &QueueItem{}) // TODO: handle error
		nextItem := resNext.(QueueItem)
		// remove PrevKey from nextItem
		nextItem.PrevKey = EmptyItemPointerKey
		// save updated nextItem
		c.State().Put(nextItem) // TODO: handle error
	}
	storeHeadKey(c, nextKey)
	if isKeyEmpty(nextKey) { // reached a TailItem
		storeTailKey(c, nextKey)
	}
	// remove extracted item from state
	c.State().Delete(headKey) // TODO: handle error

	//return item, c.State().Delete(item)
	extractedItem = headItem
	return extractedItem, nil
}

// attcahes extra data to the item specified by key, returns error if key not exists
// arg1 -> attachDataMethodParamKey string
// arg2 -> attachDataMethodParamData []bytes
func queueAttachData(c router.Context) (interface{}, error) {
	idStr := c.ParamString(itemIDParam)
	id, err := ulid.ParseStrict(idStr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ULID string passed")
	}
	extraData := c.ParamBytes(attachedDataParam)
	queryItem := &QueueItem{ID: id}
	key, _ := queryItem.Key() // always no error
	item, err := readQueueItem(c, key)
	if err != nil {
		return nil, errors.Wrap(err, "can not read item to attach data")
	}
	item.ExtraData = append(item.ExtraData, extraData...)
	// fmt.Printf("\n\n***** item=%+v\n\n", item)
	if err := c.State().Put(item); err != nil {
		return nil, errors.Wrap(err, "failed to update item with extra data")
	}
	return item, nil
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

func queueListItems(c router.Context) (interface{}, error) {
	// we can raplace realization to any of queueListItems*
	return queueListItemsItarated(c)
}

// queueListItemsMemSorted read and return all queue items as list sorted by ULID stored in ID
func queueListItemsMemSorted(c router.Context) (interface{}, error) {
	res, err := c.State().List(queueItemKeyPrefix, &QueueItem{})
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

// queueListItemsDBSorted implemens lising queue by CouchDB rich query
// NOTE: you can not test it by Mock, it doesn't implement GetQueryResult()
func queueListItemsDBSorted(c router.Context) (interface{}, error) {
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
	return c.State().List(queueItemKeyPrefix, &QueueItem{})
}

// gets a list form the ledger by travaling along Next links between nodes
func queueListItemsItarated(c router.Context) (interface{}, error) {
	// TODO: take HEAD and iterate until NEXT != *EMPTY*[
	items := []QueueItem{}
	headPresent, _ := hasHead(c) // TODO: handle error
	if !headPresent {
		return items, nil // return empty list
	}
	head, _ := getHeadItem(c) // TODO: handle error
	items = append(items, head)

	nextKey := head.NextKey
	for !isKeyEmpty(nextKey) {
		item, err := readQueueItem(c, nextKey)
		if err != nil {
			return items, errors.Wrap(err, "failed read item to list")
		}
		items = append(items, item)
		nextKey = item.NextKey
	}
	return items, nil
}
