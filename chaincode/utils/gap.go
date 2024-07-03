package utils

import (
	"encoding/json"
	"fmt"

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

	if input.IssueDateFrom != nil && input.IssueDateTo != nil {
		filter["issueDate"] = map[string]interface{}{
			"$gte": *input.IssueDateFrom,
			"$lte": *input.IssueDateTo,
		}
	}

	if input.ExpireDateFrom != nil && input.ExpireDateTo != nil {
		filter["expireDate"] = map[string]interface{}{
			"$gte": *input.ExpireDateFrom,
			"$lte": *input.ExpireDateTo,
		}
	}

	if input.AvailableGap != nil {
		filter["farmerId"] = ""
	}

	if input.Gaps != nil && len(input.Gaps) > 0 {
		filter["certId"] = map[string]interface{}{
			"$in": input.Gaps,
		}
	}

	filter["docType"] = "gap"

	filterJSON, err := json.Marshal(filter)
	if err != nil {
		fmt.Printf("Error marshalling filter to JSON: %v\n", err)
	} else {
		fmt.Printf("gap filter: %s\n", filterJSON)
	}

	return filter
}

func GapFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllGap, filter map[string]interface{}) ([]*models.GapTransactionResponse, error) {

	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip > 0 {
		selector["skip"] = input.Skip
	}
	if input.Limit > 0 {
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
