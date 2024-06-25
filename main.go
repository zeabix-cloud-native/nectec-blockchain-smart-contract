/*
SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	chaincode "github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func main() {
	smartContract, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	
	if err != nil {
		log.Panicf("Error creating nectec chaincode: %v", err)
	}

	if err := smartContract.Start(); err != nil {
		log.Panicf("Error starting nectec chaincode: %v", err)
	}
}
