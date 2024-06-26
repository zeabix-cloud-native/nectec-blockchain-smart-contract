/*
SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode"
)

type SmartContractWrapper struct {
    *chaincode.SmartContract
}

func main() {
	sc := &chaincode.SmartContract{}
    scWrapper := &SmartContractWrapper{SmartContract: sc}

    chaincode, err := contractapi.NewChaincode(scWrapper)
	
	if err != nil {
		log.Panicf("Error creating nectec chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting nectec chaincode: %v", err)
	}
}
