package hlfq

// ExampleItems is a test data
var ExampleItems = []QueueItemSpec{{
	From:      "A",
	To:        "B",
	Amount:    1,
	ExtraData: []byte("Extra Data for A:B->1"),
}, {
	From:   "B",
	To:     "C",
	Amount: 2,
}, {
	From:   "A",
	To:     "B",
	Amount: 3,
}, {
	From:      "C",
	To:        "B",
	Amount:    4,
	ExtraData: []byte("Extra Data for C:B->4"),
}}
