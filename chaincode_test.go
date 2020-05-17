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
	cc := testcc.NewMockStub("hlfqueue", hlfq.New())
	ccWithoutAC := testcc.NewMockStub("hlfqueue", hlfq.New())

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(cc.From(Authority).Init()) // init chaincode from authority
	})

	Describe("HLFQueue", func() {

		It("Allows to push an item to the queue", func() {
			//invoke chaincode method from non authority actor
			expectcc.ResponseOk(
				ccWithoutAC.From(Authority).Invoke("hlfqueuePush", hlfq.ExampleItems[0]))
		})

		It("Allows to pop an item from the queue", func() {
			//invoke chaincode method from non authority actor
			expectcc.ResponseOk(
				ccWithoutAC.From(Authority).Invoke("hlfqueuePop", hlfq.ExampleItems[0]))
		})

		It("Allow to add extra data to specified queue item", func() {
			// register second hlfqueue
			expectcc.ResponseOk(cc.From(Authority).Invoke("hlfqueueRegister", hlfq.ExampleItems[1]))
			cc := expectcc.PayloadIs(
				cc.From(Authority).Invoke("hlfqueueAttachData"),
				&[]hlfq.QueueItem{}).([]hlfq.QueueItem)

			Expect(cc).To(HaveLen(2))
		})

		It("Allow to get queue items as list", func() {
			//  &[]QueueItem{} - declares target type for unmarshalling from []byte received from chaincode
			cc := expectcc.PayloadIs(cc.Invoke("hlfqueueList"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)

			Expect(cc).To(HaveLen(1))
			Expect(cc[0].ID).To(Equal(hlfq.ExampleItems[0].ID))
		})

		It("Allows to move an item to the place AFTER specified item", func() {
			//  &[]QueueItem{} - declares target type for unmarshalling from []byte received from chaincode
			// TODO: проверять что в результате в списке в очереди элемент E, идет после After

			//cc := expectcc.PayloadIs(cc.Invoke("hlfqueueMoveAfter"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			//Expect(cc).To(HaveLen(1))
			//Expect(cc[0].ID).To(Equal(hlfq.ExampleItems[0].ID))
		})

		It("Allows to move an item to the place BEFORE specified item", func() {
			//  &[]QueueItem{} - declares target type for unmarshalling from []byte received from chaincode
			// TODO: проверять что в результате в списке в очереди элемент E, идет перед элементом Before

			//cc := expectcc.PayloadIs(cc.Invoke("hlfqueueMoveAfter"), &[]hlfq.QueueItem{}).([]hlfq.QueueItem)
			//Expect(cc).To(HaveLen(1))
			//Expect(cc[0].ID).To(Equal(hlfq.ExampleItems[0].ID))
		})
	})
})
