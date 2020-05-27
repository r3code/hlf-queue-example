# FIFO Queue Hyperledger Fabric chaincode

[![Build Status](https://travis-ci.com/r3code/hlf-queue-example.svg?branch=master)](https://travis-ci.com/r3code/hlf-queue-example)

HLFQueue chaincode stores and manages a FIFO queue.


## Supported chaincode methods

**Push** - adds an item data to the tail of the queue and returns created queue item. ID of the item generated automatically as ULID (see https://github.com/oklog/ulid).

**Pop** - dequeues (extracts) an item from the head of the queue. If queue is empty it will raise an error "Empty queue".

**Select** - allows you to filter queue items using a query string in `expr` syntax (see https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md). Returns a list of matched queue items. Example query `{.Amount > 1 and .Amount < 4}` - select items where `Amount` between 1 and 4.

**ListItems** - returns a list of all item in queue.

**Attach Data** - attaches specified `[]byte` data to an item `ExtraData` specified by `ID` (ULID string). Replaces existing item `ExtraData`.

**MoveAfter** - cuts the item and puts it after the specified item ID in the queue.

**MoveBefore** - cuts the item and puts it before the specified item ID in the queue.


## Building

### Dependencies 

See https://github.com/SAPDocuments/Tutorials/issues/4415

	go get github.com/hyperledger/fabric-chaincode-go/shim
	go get github.com/hyperledger/fabric/core/peer
	go get github.com/hyperledger/fabric/common/util


### Fix Docker package error at Windows

`C:\Users\r3code\go\pkg\mod\github.com\docker\docker@v1.4.2-0.20191101170500-ac7306503d23\pkg\system\filesys_windows.go:112:24: cannot use uintptr(unsafe.Pointer(&sd[0])) (type uintptr) as type *"golang.org/x/sys/windows".SECURITY_DESCRIPTOR in assignment`

Install a package:

	go get github.com/docker/docker@2200d938a2d5e7cd7437489c22a32d37d9bb380d


### Fix build errors

Replace in file paths a `shim`-path to 

	"github.com/hyperledger/fabric-chaincode-go/shim"

because hyperledger API path has changed.

## Build and start the chaincode

First - start network.

### Install the chaincode 

	docker exec -it chaincode bash
	// output: root@d2629980e76b:/opt/gopath/src/chaincode

**Сompile the chaincode**

	cd ./hlf-queue-example/cmd/hlfqueue
	go build

**Run the chaincode**

	CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./hlfqueue

The chaincode is started with peer and chaincode logs indicating successful registration with the peer.

**Prepare to use**

	docker exec -it cli bash
	peer chaincode install -p chaincodedev/chaincode/hlf-queue-example/cmd/hlfqueue -n mycc -v 0
	peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C mychannel

`Instantinate` will init the ledger default states used by the chaincode.


### Push an item to the queue

Push an item with data:

	{
		"from": "A",
		"to": "B",
		"amount": 1
	}

Execute a command:

	peer chaincode invoke -n mycc -c '{"Args":["Push", "{\"From\":\"A\",\"To\":\"B\", \"Amount\": 1 }"]}' -C myc

Push an item with extra data:

	{
		"From": "A",
		"To": "B",
		"Amount": 1,
		"ExtraData": "A to B"
	}

Execute a command:

	peer chaincode invoke -n mycc -c '{"Args":["Push", "{\"From\":\"A\",\"To\":\"B\", \"Amount\": 1, \"ExtraData\": \"A to B\" }"]}' -C myc

### Pop an item from the queue

	peer chaincode invoke -n mycc -c '{"Args":["Pop"]}' -C myc

### Reordering queue items

#### Move after

Cut the item with ID `01D78XYFJ1PRM1WPBCBT3VITEM` and put after `01D78XYFJ1PRM1WPBCBT3AFTER`.

	peer chaincode invoke -n mycc -c '{"Args":["MoveAfter", "01D78XYFJ1PRM1WPBCBT3VITEM", "01D78XYFJ1PRM1WPBCBT3AFTER"]}' -C myc

#### Move before

Cut the item with ID `01D78XYFJ1PRM1WPBCBT3VHOER` and put before `01D78XYFJ1PRM1WPBCBT3VHMNV`.

	peer chaincode invoke -n mycc -c '{"Args":["MoveBefore", "01D78XYFJ1PRM1WPBCBT3VHOER", "01D78XYFJ1PRM1WPBCBT3VHMNV"]}' -C myc

### Select queue items (filtering)

Select all items where `From = "A"` and `Amount > 2`

	peer chaincode invoke -n mycc -c '{"Args":["Select", "{.From == \"A\" and .Amount > 2 }"]}' -C myc

### Attach data	to an item with specified ID

	peer chaincode invoke -n mycc -c '{"Args":["AttachData", "01D78XYFJ1PRM1WPBCBT3VHMNV", "Data to attach"]}' -C myc

### Extra 

#### List queue items

	peer chaincode invoke -n mycc -c '{"Args":["ListItems"]}' -C myc


## Development

### Testing 

To run tests call:

	go test

### Debugging 

Set `CORE_CHAINCODE_LOGGING_LEVEL=debug` to see a debug output.
If you want to add your own Debug messages use `c.Logger().Debug(msg)`

In Windows PowerShell:

	$env:CORE_CHAINCODE_LOGGING_LEVEL=debug
	go test // will output "DEBU" prefixed messaged to STDOUT

In Linux:

	$CORE_CHAINCODE_LOGGING_LEVEL=debug
	go test // will output "DEBU" prefixed messaged to STDOUT



