package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func RegulatorFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllRegulator) ([]*models.RegulatorTransactionResponse, error) {
	var filter = map[string]interface{}{}

	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getStringRegulator, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryRegulator, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringRegulator), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryRegulator.Close()

	var dataRegulator []*models.RegulatorTransactionResponse
	for queryRegulator.HasNext() {
		queryRes, err := queryRegulator.Next()
		if err != nil {
			return nil, err
		}

		var dataR models.RegulatorTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataR)
		if err != nil {
			return nil, err
		}

		dataRegulator = append(dataRegulator, &dataR)
	}

	return dataRegulator, nil
}