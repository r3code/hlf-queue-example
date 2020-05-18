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
	From:   "C",
	To:     "D",
	Amount: 3,
}, {
	From:      "D",
	To:        "F",
	Amount:    4,
	ExtraData: []byte("Extra Data for D:F->4"),
}}
