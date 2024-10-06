package main

import (
	"log"

	"github.com/hyperledger/fabric-samples/chaincode/fabcar/go/contracts"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	asssetContract := new(contracts.AssetContract)
	chaincode, err := contractapi.NewChaincode(asssetContract)
	if err != nil {
		log.Fatalf("error while creating new chaincode: %v", err)
	}
	err = chaincode.Start()
	if err != nil {
		log.Fatalf("error while starting chaincode: %v", err)
	}
}
