# HLFQueue registration Hyperledger Fabric chaincode

HLFQueue registration chaincode use simple golang structure with json marshalling.



## Building

### Dependencies 

See https://github.com/SAPDocuments/Tutorials/issues/4415

go get github.com/hyperledger/fabric-chaincode-go/shim
go get github.com/hyperledger/fabric/core/peer
go get github.com/hyperledger/fabric/common/util


== Устранение ошибки Docker

@C:\Users\Дмитрий\go\pkg\mod\github.com\docker\docker@v1.4.2-0.20191101170500-ac7306503d23\pkg\system\filesys_windows.go:112:24: cannot use uintptr(unsafe.Pointer(&sd[0])) (type uintptr) as type *"golang.org/x/sys/windows".SECURITY_DESCRIPTOR in assignment@

go get github.com/docker/docker@2200d938a2d5e7cd7437489c22a32d37d9bb380d

== Устранение ошибок сборки

заменить пути к shim на 
"github.com/hyperledger/fabric-chaincode-go/shim"

