package models

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

type SmartContract struct {
	contractapi.Contract
}

type Pagination struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}