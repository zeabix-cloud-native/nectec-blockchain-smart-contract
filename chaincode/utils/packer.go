package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func PackerFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllPacker) ([]*models.PackerTransactionResponse, error) {
	var filter = map[string]interface{}{}

	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getStringPacker, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryPacker, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringPacker), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryPacker.Close()

	var dataPacker []*models.PackerTransactionResponse
	for queryPacker.HasNext() {
		queryRes, err := queryPacker.Next()
		if err != nil {
			return nil, err
		}

		var dataP models.PackerTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataP)
		if err != nil {
			return nil, err
		}

		dataPacker = append(dataPacker, &dataP)
	}

	return dataPacker, nil
}
