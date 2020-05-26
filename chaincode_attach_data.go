package hlfq

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

const attachedDataParam = "attachedData"

// attcahes extra data to the item specified by key, returns error if key not exists
// arg1 -> attachDataMethodParamKey string
// arg2 -> attachDataMethodParamData []bytes
func queueAttachData(c router.Context) (interface{}, error) {
	itemIDStr := c.ParamString(itemIDParam)
	extraData := c.ParamBytes(attachedDataParam)
	item, err := readQueueItemByID(c, itemIDStr)
	if err != nil {
		return nil, errors.Wrap(err, "can not read item to attach data")
	}
	item.ExtraData = []byte{} // reset
	item.ExtraData = append(item.ExtraData, extraData...)
	// fmt.Printf("\n\n***** item=%+v\n\n", item)
	if err := c.State().Put(item); err != nil {
		return nil, errors.Wrap(err, "failed to update item with extra data")
	}
	return item, nil
}
