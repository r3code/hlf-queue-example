# HLFQueue registration Hyperledger Fabric chaincode

HLFQueue registration chaincode use simple golang structure with json marshalling.



## Building


### Dependencies 

See https://github.com/SAPDocuments/Tutorials/issues/4415

go get github.com/hyperledger/fabric-chaincode-go/shim
go get github.com/hyperledger/fabric/core/peer
go get github.com/hyperledger/fabric/common/util


### Fix Docker package error at Windows

@C:\Users\r3code\go\pkg\mod\github.com\docker\docker@v1.4.2-0.20191101170500-ac7306503d23\pkg\system\filesys_windows.go:112:24: cannot use uintptr(unsafe.Pointer(&sd[0])) (type uintptr) as type *"golang.org/x/sys/windows".SECURITY_DESCRIPTOR in assignment@

go get github.com/docker/docker@2200d938a2d5e7cd7437489c22a32d37d9bb380d


### Fix build errors

Replace in file paths a `shim`-path to 
"github.com/hyperledger/fabric-chaincode-go/shim"
Because hyperledger api path has changed.


### Debugging 

Set `CORE_CHAINCODE_LOGGING_LEVEL=debug` to see a debug output.
If you want to add your own Debug messages use `c.Logger().Debug(msg)`

In Windows PowerShell:
	$env:CORE_CHAINCODE_LOGGING_LEVEL=debug
	go test // will output "DEBU" prefixed messaged to STDOUT

In Linux:
	$CORE_CHAINCODE_LOGGING_LEVEL=debug
	go test // will output "DEBU" prefixed messaged to STDOUT


## Supported chaincode methods

### Push

Adds an item to the tail of the queue.

### Pop

Extracts an item from the head of the queue. If queue is empty it will raise an error "Empty queue"

### Select

Allows you to filter queue items using a query string in `expr` syntax (see https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md)
Returns a list of matched queue items.

### ListItems

Returns a list of all item in queue.

### Attach Data

Attaches specified `[]byte` data to an item `ExtraData` specified by `ID` (ULID string).
Replaces existing item `ExtraData`

