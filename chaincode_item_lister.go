package hlfq

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

// ###
// ## Different ways of listing queue items
// ###

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
