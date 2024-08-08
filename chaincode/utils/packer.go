package utils

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func PackerFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllPacker, filter map[string]interface{}) ([]*models.PackerTransactionResponse, int, error) {
	search := input.Search
	filter["docType"] = "packer"

	selector := map[string]interface{}{
        "selector": filter,
		"use_index": []string{
            "_design/index-DocType",
            "index-DocType",
        },
    }

	if search != nil && *search != "" {
        searchTerm := *search
        selector["selector"] = map[string]interface{}{
            "$and": []map[string]interface{}{
                filter,
                {
                    "$or": []map[string]interface{}{
                        {"packingHouseRegisterNumber": map[string]interface{}{"$regex": searchTerm}},
                        {"packingHouseName": map[string]interface{}{"$regex": searchTerm}},
                    },
                },
            },
        }
    }

	getStringPacker, err := json.Marshal(selector)
	if err != nil {
		return nil, 0, err
	}

	// Fetch total count of the results
	resultsIterator, err := ctx.GetStub().GetQueryResult(string(getStringPacker))
	if err != nil {
		return nil, 0, err
	}
	defer resultsIterator.Close()

	total := 0
	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return nil, 0, err
		}
		total++
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	fmt.Printf("packer filter %v", selector)

	queryPacker, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringPacker), int32(input.Limit), "")
	if err != nil {
		return nil, 0, err
	}
	defer queryPacker.Close()

	var dataPacker []*models.PackerTransactionResponse
	for queryPacker.HasNext() {
		queryRes, err := queryPacker.Next()
		if err != nil {
			return nil, 0, err
		}

		var dataP models.PackerTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataP)
		if err != nil {
			return nil, 0, err
		}

		dataPacker = append(dataPacker, &dataP)
	}

	return dataPacker, total, nil
}
