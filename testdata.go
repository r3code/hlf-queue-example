package hlfq

// ExampleItems is a test data
var ExampleItems = []*QueueItemSpec{{
	ID:        "1",
	From:      "A",
	To:        "B",
	Amount:    1,
	ExtraData: []byte("Extra Data for 1"),
}, {
	ID:     "2",
	From:   "A",
	To:     "B",
	Amount: 4,
}, {
	ID:     "3",
	From:   "A",
	To:     "C",
	Amount: 6,
}, {
	ID:     "4",
	From:   "B",
	To:     "C",
	Amount: 8,
}}
