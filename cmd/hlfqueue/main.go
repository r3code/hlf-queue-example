package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	hlfq "github.com/r3code/hlf-queue-example"
)

func main() {
	cc := hlfq.New()
	if err := shim.Start(cc); err != nil {
		fmt.Printf("Error starting HLFQueue chaincode: %s", err)
		os.Exit(1)
	}

}
