package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func NectecStaffFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllNectecStaff, filter map[string]interface{}) ([]*models.NectecStaffTransactionResponse, int, error) {
	filter["docType"] = "nectec"

	selector := map[string]interface{}{
		"selector": filter,
	}

	getStringNectec, err := json.Marshal(selector)
    if err != nil {
        return nil, 0, err
    }

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(getStringNectec))
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

	if input.Skip > 0 {
        selector["skip"] = input.Skip
    }
    if input.Limit > 0 {
        selector["limit"] = input.Limit
    }

	getStringNstda, err := json.Marshal(selector)
	if err != nil {
		return nil, 0, err
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

	return dataNstda, total, nil
}


