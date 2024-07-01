package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func GmpSetFilter(input *models.FilterGetAllGmp) map[string]interface{} {
	var filter = map[string]interface{}{}

	if input.PackingHouseRegisterNumber != nil {
		filter["packingHouseRegisterNumber"] = input.PackingHouseRegisterNumber
	}

	if input.PackingHouseName != nil {
		filter["packingHouseName"] = map[string]interface{}{
			"$regex": input.PackingHouseName,
		}
	}

	if input.Address != nil {
		filter["address"] = map[string]interface{}{
			"$regex": input.Address,
		}
	}

	filter["docType"] = "gmp"

	return filter
}

func GmpFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllGmp, filter map[string]interface{}) ([]*models.GmpTransactionResponse, error) {
	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getStringGmp, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryGmp, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringGmp), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryGmp.Close()

	var dataGmp []*models.GmpTransactionResponse
	for queryGmp.HasNext() {
		queryRes, err := queryGmp.Next()
		if err != nil {
			return nil, err
		}

		var dataG models.GmpTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataG)
		if err != nil {
			return nil, err
		}

		dataGmp = append(dataGmp, &dataG)
	}

	return dataGmp, nil
}