package utils

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func FarmerFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllFarmer) ([]*models.FarmerTransactionResponse, int, error) {
	var filter = map[string]interface{}{}

	filter["docType"] = "farmer"

	if input.Search != "" {
		filter["$or"] = []map[string]interface{}{
			{
				"farmerGaps": map[string]interface{}{
					"$elemMatch": map[string]interface{}{
						"displayCertId": map[string]interface{}{ 
							"$regex": input.Search,
						},
					},
				},
			},
			{
				"id": input.Search,
			},
		}
	}

	selector := map[string]interface{}{
		"selector": filter,
		"sort": []map[string]string{
			{"createdAt": "desc"},
		},
		"use_index": []string{
			"_design/indexCreatedAtId",
			"indexCreatedAtId",
		},
	}

	if input.Skip > 0 {
		selector["skip"] = input.Skip
	}
	if input.Limit > 0 {
		selector["limit"] = input.Limit
	}

	getString, err := json.Marshal(selector)
	if err != nil {
		fmt.Printf("Error marshalling filter to JSON: %v\n", err)
	} else {
		fmt.Printf("farmer filter: %s\n", getString)
	}

	queryFarmer, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getString), int32(input.Limit), "")
	if err != nil {
		return nil, 0, err
	}
	defer queryFarmer.Close()

	var dataFarmers []*models.FarmerTransactionResponse
	for queryFarmer.HasNext() {
		queryRes, err := queryFarmer.Next()
		if err != nil {
			return nil, 0, err
		}

		var dataF models.FarmerTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataF)
		if err != nil {
			return nil, 0, err
		}

		// Ensure FarmerGaps is not nil
		if dataF.FarmerGaps == nil {
			dataF.FarmerGaps = []models.FarmerGap{}
		}

		dataFarmers = append(dataFarmers, &dataF)
	}

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(getString))
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

	return dataFarmers, total, nil
}
