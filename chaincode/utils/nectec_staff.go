package utils

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func NectecStaffFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllNectecStaff, filter map[string]interface{}) ([]*models.NectecStaffTransactionResponse, int, error) {
	filter["docType"] = "nectec"

	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Search != nil && *input.Search != "" {
		filter["certId"] = map[string]interface{}{
			"$regex": input.Search,
		}
	}

	selector = map[string]interface{}{
		"selector": filter,
		"sort": []map[string]string{
			{"createdAt": "desc"},
		},
		"use_index": []string{
			"_design/indexCreatedAtId",
			"indexCreatedAtId",
		},
	}

	getStringNstda, err := json.Marshal(selector)
	if err != nil {
		return nil, 0, err
	} else {
		fmt.Printf("staff query filter: %s\n", getStringNstda)
	}

	total, err := CountTotalResults(ctx, string(getStringNstda))
	if err != nil {
		return nil, 0, err
	}

	if input.Skip > 0 {
        selector["skip"] = input.Skip
    }
    if input.Limit > 0 {
        selector["limit"] = input.Limit
    }

	getQueryString, err := json.Marshal(selector)
	if err != nil {
		return nil, 0, err
	} else {
		fmt.Printf("staff query filter: %s\n", getQueryString)
	}
	queryNstda, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringNstda), int32(input.Limit), "")
	if err != nil {
		return nil, 0, err
	}
	defer queryNstda.Close()

	var dataNstda []*models.NectecStaffTransactionResponse
	for queryNstda.HasNext() {
		queryRes, err := queryNstda.Next()
		if err != nil {
			return nil, 0, err
		}

		var dataP models.NectecStaffTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataP)
		if err != nil {
			return nil, 0, err
		}

		dataNstda = append(dataNstda, &dataP)
	}

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(getStringNstda))
    if err != nil {
        return nil, 0, err
    }
    defer resultsIterator.Close()

	return dataNstda, total, nil
}


