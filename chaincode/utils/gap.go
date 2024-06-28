package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func GapSetFilter(input *models.FilterGetAllGap) map[string]interface{} {
	var filter = map[string]interface{}{}
	
	if input.FarmerID != nil {
		filter["farmerId"] = *input.FarmerID
	}
	if input.CertID != nil {
		filter["certId"] = *input.CertID
	}
	if input.AreaCode != nil {
		filter["areaCode"] = *input.AreaCode
	}
	if input.Province != nil {
		filter["province"] = *input.Province
	}
	if input.District != nil {
		filter["district"] = *input.District
	}
	if input.AreaRaiFrom != nil && input.AreaRaiTo != nil {
		filter["areaRai"] = map[string]interface{}{
			"$gte": *input.AreaRaiFrom,
			"$lte": *input.AreaRaiTo,
		}
	}
	if input.IssueDate != nil {
		filter["issueDate"] = *input.IssueDate
	}
	if input.ExpireDate != nil {
		filter["expireDate"] = *input.ExpireDate
	}

	if input.AvailableGap != nil {
		filter["farmerId"] = ""
	}

	filter["docType"] = "gap"

	return filter
}

func GapFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllGap, filter map[string]interface{}) ([]*models.GapTransactionResponse, error) {

	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	queryString, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryResults, _, err := ctx.GetStub().GetQueryResultWithPagination(string(queryString), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryResults.Close()

	var assets []*models.GapTransactionResponse
	for queryResults.HasNext() {
		queryResponse, err := queryResults.Next()
		if err != nil {
			return nil, err
		}

		var asset models.GapTransactionResponse
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	return assets, nil
}
