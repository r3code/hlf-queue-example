package hlfq_test

import (
	"testing"

	hlfq "github.com/r3code/hlf-queue-example"
	"github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHLFQueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HLFQueue Suite")
}

var (
	Authority = testdata.Certificates[0].MustIdentity("SOME_MSP")
	Someone   = testdata.Certificates[1].MustIdentity("SOME_MSP")
)

var _ = Describe("HLFQueue", func() {

	//Create chaincode mock
	ccMockGlobal := testcc.NewMockStub("hlfq_mock", hlfq.New())

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(ccMockGlobal.From(Authority).Init()) // init chaincode
	})

	// Describe("Push/Pop", func() {

	// 	It("Allows to push an item to the queue", func() {
	// 		testData := hlfq.ExampleItems[0]
	// 		expectcc.ResponseOk(
	// 			ccMockGlobal.From(Authority).Invoke("Push", testData))
	// 		// get list and check it has one expected element
	// 		items := expectcc.PayloadIs(ccMockGlobal.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
	// 		Expect(items).To(HaveLen(1))
	// 		Expect(items[0].From).To(Equal(testData.From))
	// 		Expect(items[0].To).To(Equal(testData.To))
	// 		Expect(items[0].Amount).To(Equal(testData.Amount))
	// 	})

	// 	It("Allows to pop an item from the queue", func() {
	// 		//invoke chaincode method from non authority actor
	// 		headitem := expectcc.PayloadIs(ccMockGlobal.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
	// 		Expect(headitem.From).To(Equal(hlfq.ExampleItems[0].From))
	// 		Expect(headitem.To).To(Equal(hlfq.ExampleItems[0].To))
	// 		Expect(headitem.Amount).To(Equal(hlfq.ExampleItems[0].Amount))

	// 		// get list and check it has 0 items now
	// 		items := expectcc.PayloadIs(ccMockGlobal.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
	// 		Expect(items).To(HaveLen(0))
	// 	})

	// 	It("Push 3 items and pops 3 items and checks if it in FIFO order", func() {
	// 		//invoke chaincode method from non authority actor
	// 		// Push 3 items
	// 		expectcc.ResponseOk(
	// 			ccMockGlobal.From(Authority).Invoke("Push", hlfq.ExampleItems[0]))
	// 		expectcc.ResponseOk(
	// 			ccMockGlobal.From(Authority).Invoke("Push", hlfq.ExampleItems[1]))
	// 		expectcc.ResponseOk(
	// 			ccMockGlobal.From(Authority).Invoke("Push", hlfq.ExampleItems[2]))
	// 		headItem1 := expectcc.PayloadIs(ccMockGlobal.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
	// 		Expect(headItem1.From).To(Equal(hlfq.ExampleItems[0].From))
	// 		Expect(headItem1.To).To(Equal(hlfq.ExampleItems[0].To))
	// 		Expect(headItem1.Amount).To(Equal(hlfq.ExampleItems[0].Amount))
	// 		//
	// 		headItem2 := expectcc.PayloadIs(ccMockGlobal.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
	// 		Expect(headItem2.From).To(Equal(hlfq.ExampleItems[1].From))
	// 		Expect(headItem2.To).To(Equal(hlfq.ExampleItems[1].To))
	// 		Expect(headItem2.Amount).To(Equal(hlfq.ExampleItems[1].Amount))
	// 		//
	// 		headItem3 := expectcc.PayloadIs(ccMockGlobal.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
	// 		Expect(headItem3.From).To(Equal(hlfq.ExampleItems[2].From))
	// 		Expect(headItem3.To).To(Equal(hlfq.ExampleItems[2].To))
	// 		Expect(headItem3.Amount).To(Equal(hlfq.ExampleItems[2].Amount))

	// 		// get list and check it has 0 items now
	// 		items := expectcc.PayloadIs(ccMockGlobal.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
	// 		Expect(items).To(HaveLen(0))
	// 	})

	// })

	// Describe("Inspect Queue", func() {

	// 	It("Allow to get queue items as list", func() {
	// 		// Add one item
	// 		testItems := hlfq.ExampleItems[0:3]
	// 		for _, ti := range testItems {
	// 			ccMockGlobal.From(Authority).Invoke("Push", ti)
	// 			// NOTE: if there is no sleep you would get an unexpected order of elemenst in List,
	// 			time.Sleep(time.Millisecond) // if no delay, Push() #3 of #2 can be executed before Push #1
	// 		}
	// 		//  &[]QueueItem{} - declares target type for unmarshalling from []byte received from chaincode
	// 		items := expectcc.PayloadIs(ccMockGlobal.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
	// 		// fmt.Printf("  ***TEST_items=%+v\n", testItems)
	// 		// fmt.Printf("   ***     items=%+v\n", items)
	// 		Expect(items).To(HaveLen(3))

	// 		for i, ti := range testItems {
	// 			// Dont comapre ID, it's generated at any Push call
	// 			Expect(items[i].Amount).To(Equal(ti.Amount))
	// 			Expect(items[i].From).To(Equal(ti.From))
	// 			Expect(items[i].To).To(Equal(ti.To))

	// 		}

	// 	})
	// })

	// Describe("Attach extra context to item", func() {

	// 	It("Allow to add extra data to specified queue item", func() {
	// 		// lets Push 3 items into queue in reverse order
	// 		expectcc.ResponseOk(
	// 			ccMockGlobal.From(Authority).Invoke("Push", hlfq.ExampleItems[2])) // head
	// 		expectcc.ResponseOk(
	// 			ccMockGlobal.From(Authority).Invoke("Push", hlfq.ExampleItems[1]))
	// 		// at the begin the test item has no ExtraData
	// 		Expect(hlfq.ExampleItems[1].ExtraData).To(Equal([]byte{}))
	// 		expectcc.ResponseOk(
	// 			ccMockGlobal.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // tail
	// 		// take a list
	// 		items := expectcc.PayloadIs(ccMockGlobal.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
	// 		// select item 1
	// 		item1 := items[1]
	// 		item1IDStr := item1.ID.String()
	// 		// attach extra data
	// 		testExtraData := []byte("An extra data for " + item1IDStr)

	// 		updatedItem := expectcc.PayloadIs(
	// 			ccMockGlobal.From(Authority).Invoke("AttachData", item1IDStr, testExtraData),
	// 			&hlfq.QueueItem{}).(hlfq.QueueItem)

	// 		Expect(updatedItem.ID).To(Equal(item1.ID))
	// 		Expect(string(updatedItem.ExtraData)).To(Equal(string(testExtraData)))
	// 	})

	// })

	// Describe("Check Select by Filter", func() {

	// 	It("Selects items with specified Amount range", func() {
	// 		// we need new empty queue
	// 		ccMock2 := testcc.NewMockStub("hlfq_mock2", hlfq.New())
	// 		expectcc.ResponseOk(ccMock2.From(Authority).Init()) // init chaincode2

	// 		// Push 3 items
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // Amount=1
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[1])) // Amount=2
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[2])) // Amount=3
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[3])) // Amount=4

	// 		queryStr := "{.Amount > 1 and .Amount < 4}"
	// 		filteredItems := expectcc.PayloadIs(
	// 			ccMock2.From(Authority).Invoke("Select", queryStr),
	// 			&[]hlfq.QueueItem{}).([]hlfq.QueueItem)

	// 		Expect(filteredItems).To(HaveLen(2))
	// 		Expect(filteredItems[0].Amount).To(Equal(hlfq.ExampleItems[1].Amount))
	// 		Expect(filteredItems[1].Amount).To(Equal(hlfq.ExampleItems[2].Amount))
	// 	})

	// 	It("Should select 0 items when Amount is out of range", func() {
	// 		// we need new empty queue
	// 		ccMock2 := testcc.NewMockStub("hlfq_mock2", hlfq.New())
	// 		expectcc.ResponseOk(ccMock2.From(Authority).Init()) // init chaincode2

	// 		// Push 3 items
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // Amount=1
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[1])) // Amount=2
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[2])) // Amount=3
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[3])) // Amount=4

	// 		queryStr := "{.Amount > 100 }"
	// 		filteredItems := expectcc.PayloadIs(
	// 			ccMock2.From(Authority).Invoke("Select", queryStr),
	// 			&[]hlfq.QueueItem{}).([]hlfq.QueueItem)

	// 		Expect(filteredItems).To(HaveLen(0))
	// 	})

	// 	It("Should select 2 items where From = 'A'", func() {
	// 		// we need new empty queue
	// 		ccMock2 := testcc.NewMockStub("hlfq_mock2", hlfq.New())
	// 		expectcc.ResponseOk(ccMock2.From(Authority).Init()) // init chaincode2

	// 		// Push 3 items
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // Amount=1
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[1])) // Amount=2
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[2])) // Amount=3
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[3])) // Amount=4

	// 		From := "A"
	// 		queryStr := fmt.Sprintf("{.From == '%s' }", From)
	// 		filteredItems := expectcc.PayloadIs(
	// 			ccMock2.From(Authority).Invoke("Select", queryStr),
	// 			&[]hlfq.QueueItem{}).([]hlfq.QueueItem)

	// 		Expect(filteredItems).To(HaveLen(2))
	// 		Expect(filteredItems[0].From).To(Equal(From))
	// 		Expect(filteredItems[0].Amount).To(Equal(hlfq.ExampleItems[0].Amount))
	// 		Expect(filteredItems[1].From).To(Equal(From))
	// 		Expect(filteredItems[1].Amount).To(Equal(hlfq.ExampleItems[2].Amount))
	// 	})

	// 	It("Should select 1 item where From = 'A' and Amount > 2", func() {
	// 		// we need new empty queue
	// 		ccMock2 := testcc.NewMockStub("hlfq_mock2", hlfq.New())
	// 		expectcc.ResponseOk(ccMock2.From(Authority).Init()) // init chaincode2

	// 		// Push 3 items
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // Amount=1
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[1])) // Amount=2
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[2])) // Amount=3
	// 		expectcc.ResponseOk(
	// 			ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[3])) // Amount=4

	// 		From := "A"
	// 		Amount := 2
	// 		queryStr := fmt.Sprintf("{.From == '%s' and .Amount > %d }", From, Amount)
	// 		filteredItems := expectcc.PayloadIs(
	// 			ccMock2.From(Authority).Invoke("Select", queryStr),
	// 			&[]hlfq.QueueItem{}).([]hlfq.QueueItem)

	// 		Expect(filteredItems).To(HaveLen(1))
	// 		Expect(filteredItems[0].From).To(Equal(From))
	// 		Expect(filteredItems[0].Amount).To(BeNumerically(">=", Amount))
	// 	})

	// })

	Describe("Items Rrordering :: MoveAfter", func() {

		It("Allows to move an item to the place AFTER specified item in the middle", func() {
			// we need new empty queue
			ccMock3 := testcc.NewMockStub("hlfq_mock3", hlfq.New())
			expectcc.ResponseOk(ccMock3.From(Authority).Init()) // init chaincode2

			// Push 3 items
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // Amount=1
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[1])) // Amount=2
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[2])) // Amount=3
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[3])) // Amount=4
			// We expect to see Amout list: 1, 3, 2, 4
			itemsInQueue := expectcc.PayloadIs(ccMock3.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(itemsInQueue).To(HaveLen(4))

			movingItem := itemsInQueue[1] // Amount = 2
			afterItem := itemsInQueue[2]  // Amount = 3

			movedItem := expectcc.PayloadIs(
				ccMock3.From(Authority).Invoke("MoveAfter", movingItem.ID.String(), afterItem.ID.String()),
				&hlfq.QueueItem{}).(hlfq.QueueItem)

			// check method returned the same item that was passed in
			Expect(movedItem.From).To(Equal(movingItem.From), "")
			Expect(movedItem.To).To(Equal(movingItem.To))
			Expect(movedItem.ID.String()).To(Equal(movingItem.ID.String()))

			// check list is reordered
			reorderedList := expectcc.PayloadIs(ccMock3.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(reorderedList).To(HaveLen(len(itemsInQueue)))

			Expect(reorderedList[0].Amount).To(Equal(itemsInQueue[0].Amount), "#1 ID"+reorderedList[0].ID.String())
			Expect(reorderedList[1].Amount).To(Equal(itemsInQueue[2].Amount), "#2 ID"+reorderedList[1].ID.String())
			Expect(reorderedList[2].Amount).To(Equal(itemsInQueue[1].Amount), "#3 ID"+reorderedList[2].ID.String())
			Expect(reorderedList[3].Amount).To(Equal(itemsInQueue[3].Amount), "#4 ID"+reorderedList[3].ID.String())
		})

		It("Allows to move a head item to the place AFTER tail", func() {
			// we need new empty queue
			ccMock3 := testcc.NewMockStub("hlfq_mock3", hlfq.New())
			expectcc.ResponseOk(ccMock3.From(Authority).Init()) // init chaincode2

			// Push 3 items
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // Amount=1
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[1])) // Amount=2
			// We expect to see Amout list: 1, 3, 2, 4
			itemsInQueue := expectcc.PayloadIs(ccMock3.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(itemsInQueue).To(HaveLen(2), "Added 2 items to queue")

			movingItem := itemsInQueue[0] // Amount = 2
			afterItem := itemsInQueue[1]  // Amount = 3

			movedItem := expectcc.PayloadIs(
				ccMock3.From(Authority).Invoke("MoveAfter", movingItem.ID.String(), afterItem.ID.String()),
				&hlfq.QueueItem{}).(hlfq.QueueItem)

			// check method returned the same item that was passed in
			Expect(movedItem.From).To(Equal(movingItem.From), "")
			Expect(movedItem.To).To(Equal(movingItem.To))
			Expect(movedItem.ID.String()).To(Equal(movingItem.ID.String()))

			// check list is reordered
			reorderedList := expectcc.PayloadIs(ccMock3.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(reorderedList).To(HaveLen(len(itemsInQueue)), "Want same items conunt in queue after reorder")

			Expect(reorderedList[0].Amount).To(Equal(itemsInQueue[1].Amount), "#1 ID"+reorderedList[0].ID.String())
			Expect(reorderedList[1].Amount).To(Equal(itemsInQueue[0].Amount), "#2 ID"+reorderedList[1].ID.String())
		})

		It("Allows to move a Tail item to the place AFTER head", func() {
			// we need new empty queue
			ccMock3 := testcc.NewMockStub("hlfq_mock3", hlfq.New())
			expectcc.ResponseOk(ccMock3.From(Authority).Init()) // init chaincode2

			// Push 3 items
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // Amount=1
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[1])) // Amount=2
			expectcc.ResponseOk(
				ccMock3.From(Authority).Invoke("Push", hlfq.ExampleItems[2])) // Amount=3
			// We expect to see Amout list: 1, 3, 2, 4
			itemsInQueue := expectcc.PayloadIs(ccMock3.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(itemsInQueue).To(HaveLen(3), "Added 3 items to queue")

			movingItem := itemsInQueue[0] // Amount = 2
			afterItem := itemsInQueue[2]  // Amount = 3

			movedItem := expectcc.PayloadIs(
				ccMock3.From(Authority).Invoke("MoveAfter", movingItem.ID.String(), afterItem.ID.String()),
				&hlfq.QueueItem{}).(hlfq.QueueItem)

			// check method returned the same item that was passed in
			Expect(movedItem.From).To(Equal(movingItem.From), "")
			Expect(movedItem.To).To(Equal(movingItem.To))
			Expect(movedItem.ID.String()).To(Equal(movingItem.ID.String()))

			// check list is reordered
			reorderedList := expectcc.PayloadIs(ccMock3.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(reorderedList).To(HaveLen(len(itemsInQueue)), "Want same items conunt in queue after reorder")

			Expect(reorderedList[0].Amount).To(Equal(itemsInQueue[1].Amount), "#1 ID"+reorderedList[0].ID.String())
			Expect(reorderedList[1].Amount).To(Equal(itemsInQueue[2].Amount), "#2 ID"+reorderedList[1].ID.String())
			Expect(reorderedList[2].Amount).To(Equal(itemsInQueue[0].Amount), "#3 ID"+reorderedList[2].ID.String())
		})
	})

	Describe("Items Rrordering :: MoveBefore", func() {
		// It("Allows to move an item to the place BEFORE specified item", func() {
		// 	// we need new empty queue
		// 	ccMock2 := testcc.NewMockStub("hlfq_mock2", hlfq.New())
		// 	expectcc.ResponseOk(ccMock2.From(Authority).Init()) // init chaincode2

		// 	// Push 3 items
		// 	expectcc.ResponseOk(
		// 		ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[0])) // Amount=1
		// 	expectcc.ResponseOk(
		// 		ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[1])) // Amount=2
		// 	expectcc.ResponseOk(
		// 		ccMock2.From(Authority).Invoke("Push", hlfq.ExampleItems[2])) // Amount=3

		// 	itemsInQueue := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)

		// 	// put item[2] before item[1]
		// 	movedItem := expectcc.PayloadIs(
		// 		ccMock2.From(Authority).Invoke("MoveBefore", itemsInQueue[2].ID.String(), itemsInQueue[1].ID.String()),
		// 		&hlfq.QueueItem{}).(hlfq.QueueItem)

		// 	// check method returned the same item that was passed in
		// 	Expect(movedItem.From).To(Equal(itemsInQueue[0].From))
		// 	Expect(movedItem.To).To(Equal(itemsInQueue[0].To))
		// 	Expect(movedItem.ID.String()).To(Equal(itemsInQueue[0].ID.String()))
		// 	// check list is reordered
		// 	reorderedList := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
		// 	Expect(reorderedList[0].ID.String()).To(Equal(itemsInQueue[0].ID.String()))
		// 	Expect(reorderedList[1].ID.String()).To(Equal(itemsInQueue[2].ID.String()))
		// 	Expect(reorderedList[2].ID.String()).To(Equal(itemsInQueue[1].ID.String()))
		// })
	})

})
