package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func NectecStaffFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllNectecStaff) ([]*models.NectecStaffTransactionResponse, error) {
	var filter = map[string]interface{}{}

	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getStringNstda, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryNstda, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringNstda), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryNstda.Close()

	var dataNstda []*models.NectecStaffTransactionResponse
	for queryNstda.HasNext() {
		queryRes, err := queryNstda.Next()
		if err != nil {
			return nil, err
		}

		var dataP models.NectecStaffTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataP)
		if err != nil {
			return nil, err
		}

		dataNstda = append(dataNstda, &dataP)
	}

	return dataNstda, nil
}