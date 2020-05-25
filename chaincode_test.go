package hlfq_test

import (
	"testing"
	"time"

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

	Describe("Push/Pop", func() {

		It("Allows to push an item to the queue", func() {
			testData := hlfq.ExampleItems[0]
			expectcc.ResponseOk(
				ccMock.From(Authority).Invoke("Push", testData))
			// get list and check it has one expected element
			items := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(items).To(HaveLen(1))
			Expect(items[0].From).To(Equal(testData.From))
			Expect(items[0].To).To(Equal(testData.To))
			Expect(items[0].Amount).To(Equal(testData.Amount))
		})

		It("Allows to pop an item from the queue", func() {
			//invoke chaincode method from non authority actor
			headitem := expectcc.PayloadIs(ccMock.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
			Expect(headitem.From).To(Equal(hlfq.ExampleItems[0].From))
			Expect(headitem.To).To(Equal(hlfq.ExampleItems[0].To))
			Expect(headitem.Amount).To(Equal(hlfq.ExampleItems[0].Amount))

			// get list and check it has 0 items now
			items := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(items).To(HaveLen(0))
		})

		It("Push 3 items and pops 3 items and checks if it in FIFO order", func() {
			//invoke chaincode method from non authority actor
			// Push 3 items
			expectcc.ResponseOk(
				ccMock.From(Authority).Invoke("Push", hlfq.ExampleItems[0]))
			expectcc.ResponseOk(
				ccMock.From(Authority).Invoke("Push", hlfq.ExampleItems[1]))
			expectcc.ResponseOk(
				ccMock.From(Authority).Invoke("Push", hlfq.ExampleItems[2]))
			headItem1 := expectcc.PayloadIs(ccMock.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
			Expect(headItem1.From).To(Equal(hlfq.ExampleItems[0].From))
			Expect(headItem1.To).To(Equal(hlfq.ExampleItems[0].To))
			Expect(headItem1.Amount).To(Equal(hlfq.ExampleItems[0].Amount))
			//
			headItem2 := expectcc.PayloadIs(ccMock.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
			Expect(headItem2.From).To(Equal(hlfq.ExampleItems[1].From))
			Expect(headItem2.To).To(Equal(hlfq.ExampleItems[1].To))
			Expect(headItem2.Amount).To(Equal(hlfq.ExampleItems[1].Amount))
			//
			headItem3 := expectcc.PayloadIs(ccMock.Invoke("Pop"), &hlfq.QueueItem{}).(hlfq.QueueItem)
			Expect(headItem3.From).To(Equal(hlfq.ExampleItems[2].From))
			Expect(headItem3.To).To(Equal(hlfq.ExampleItems[2].To))
			Expect(headItem3.Amount).To(Equal(hlfq.ExampleItems[2].Amount))

			// get list and check it has 0 items now
			items := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			Expect(items).To(HaveLen(0))
		})

	})

	Describe("Inspect Queue", func() {

		It("Allow to get queue items as list", func() {
			// Add one item
			testItems := hlfq.ExampleItems[0:3]
			for _, ti := range testItems {
				ccMock.From(Authority).Invoke("Push", ti)
				// NOTE: if there is no sleep you would get an unexpected order of elemenst in List,
				time.Sleep(time.Millisecond) // if no delay, Push() #3 of #2 can be executed before Push #1
			}
			//  &[]QueueItem{} - declares target type for unmarshalling from []byte received from chaincode
			items := expectcc.PayloadIs(ccMock.Invoke("ListItems"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			// fmt.Printf("  ***TEST_items=%+v\n", testItems)
			// fmt.Printf("   ***     items=%+v\n", items)
			Expect(items).To(HaveLen(3))

			for i, ti := range testItems {
				// Dont comapre ID, it's generated at any Push call
				Expect(items[i].Amount).To(Equal(ti.Amount))
				Expect(items[i].From).To(Equal(ti.From))
				Expect(items[i].To).To(Equal(ti.To))

			}

		})
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
