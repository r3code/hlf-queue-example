package hlfq

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
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
