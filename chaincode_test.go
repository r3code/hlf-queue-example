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
	ccMock := testcc.NewMockStub("hlfq_mock", hlfq.New())

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(ccMock.From(Authority).Init()) // init chaincode
	})

	Describe("Inspect Queue", func() {

		It("Allow to get queue items as list", func() {
			// Add one item
			testItem1 := hlfq.ExampleItems[0]
			testItem2 := hlfq.ExampleItems[1]
			ccMock.From(Authority).Invoke("Push", testItem1)
			ccMock.From(Authority).Invoke("Push", testItem2)
			//  &[]QueueItem{} - declares target type for unmarshalling from []byte received from chaincode

			items := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)

			Expect(items).To(HaveLen(2))
			// fmt.Printf("items=%v", items)
			Expect(items[0].From).To(Equal(testItem1.From))
			Expect(items[0].To).To(Equal(testItem1.To))
			Expect(items[0].Amount).To(Equal(testItem1.Amount))
			Expect(items[1].From).To(Equal(testItem2.From))
			Expect(items[1].To).To(Equal(testItem2.To))
			Expect(items[1].Amount).To(Equal(testItem2.Amount))
		})
	})

	Describe("Push/Pop", func() {

		It("Allows to push an item to the queue", func() {
			testData := hlfq.ExampleItems[3]
			expectcc.ResponseOk(
				ccMock.From(Authority).Invoke("Push", testData))
			// get list and check it has one expected element
			items := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(items).To(HaveLen(3))
			Expect(items[2].From).To(Equal(testData.From))
			Expect(items[2].To).To(Equal(testData.To))
			Expect(items[2].Amount).To(Equal(testData.Amount))
		})

		// It("Allows to pop an item from the queue", func() {
		// 	//invoke chaincode method from non authority actor
		// 	item := expectcc.PayloadIs(ccMock.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
		// 	Expect(item.From).To(Equal(hlfq.ExampleItems[2].From))
		// 	Expect(item.To).To(Equal(hlfq.ExampleItems[2].To))
		// 	Expect(item.Amount).To(Equal(hlfq.ExampleItems[2].Amount))

		// 	// get list and check it has 2 items now
		// 	items := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
		// 	Expect(items).To(HaveLen(2))
		// })

	})

	// Describe("Attach extra context to item", func() {

	// 	It("Allow to add extra data to specified queue item", func() {
	// 		// register second hlfqueue
	// 		expectcc.ResponseOk(ccMock.From(Authority).Invoke("AttachData", hlfq.ExampleItems[1]))
	// 		cc := expectcc.PayloadIs(
	// 			ccMock.From(Authority).Invoke("AttachData"),
	// 			&[]hlfq.QueueItem{}).([]hlfq.QueueItem)

	// 		Expect(cc).To(HaveLen(2))
	// 	})

	// })

	// Describe("Items Rrordering", func() {

	// 	It("Allows to move an item to the place AFTER specified item", func() {
	// 		//  &[]QueueItem{} - declares target type for unmarshalling from []byte received from chaincode
	// 		// TODO: проверять что в результате в списке в очереди элемент E, идет после After

	// 		//cc := expectcc.PayloadIs(cc.Invoke("hlfqueueMoveAfter"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
	// 		//Expect(cc).To(HaveLen(1))
	// 		//Expect(cc[0].ID).To(Equal(hlfq.ExampleItems[0].ID))
	// 	})

	// 	It("Allows to move an item to the place BEFORE specified item", func() {
	// 		//  &[]QueueItem{} - declares target type for unmarshalling from []byte received from chaincode
	// 		// TODO: проверять что в результате в списке в очереди элемент E, идет перед элементом Before

	// 		//cc := expectcc.PayloadIs(cc.Invoke("hlfqueueMoveAfter"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
	// 		//Expect(cc).To(HaveLen(1))
	// 		//Expect(cc[0].ID).To(Equal(hlfq.ExampleItems[0].ID))
	// 	})
	// })

})
