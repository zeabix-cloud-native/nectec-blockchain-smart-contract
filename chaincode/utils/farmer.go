package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func FarmerFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllFarmer) ([]*models.FarmerTransactionResponse, error) {
	var filter = map[string]interface{}{}

	filter["docType"] = "farmer"

	if input.FarmerGap != "" {
		filter["farmerGaps"] = map[string]interface{}{
			"$elemMatch": map[string]interface{}{
				"certId": input.FarmerGap,
			},
		}
	}
	
	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getString, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryFarmer, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getString), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryFarmer.Close()

	var dataFarmers []*models.FarmerTransactionResponse
	for queryFarmer.HasNext() {
		queryRes, err := queryFarmer.Next()
		if err != nil {
			return nil, err
		}

		var dataF models.FarmerTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataF)
		if err != nil {
			return nil, err
		}
		
		 if dataF.FarmerGap == nil {
			dataF.FarmerGap = []models.FarmerGap{}
	}

		dataFarmers = append(dataFarmers, &dataF)
	}

	return dataFarmers, nil
}