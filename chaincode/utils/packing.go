package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func PackingSetFilter(input *models.FilterGetAllPacking) map[string]interface{} {
	var filter = map[string]interface{}{}

	const RegexKey = "$regex"

	if input.Gap != nil {
		filter["gap"] = *input.Gap
	}

	if input.CertID != nil {
		filter["certId"] = *input.CertID
	}

	if input.Search != nil {
		filter["search"] = *input.Search
	}

	if (input.FarmerID != nil) {
		filter["farmerId"] = *input.FarmerID
	}

	if (input.PackerId != nil) {
		filter["packerId"] = *input.PackerId
	}
	
	if input.Province != nil {
		filter["province"] = *input.Province
	}
	if input.District != nil {
		filter["district"] = *input.District
	}

	if input.StartDate != nil && input.EndDate != nil {
		fromDate, err1 := FormatDate(*input.StartDate, false, offset)
		toDate, err2 := FormatDate(*input.EndDate, true, offset)

		if err1 == nil && err2 == nil {
			filter["createdAt"] = map[string]interface{}{
				"$gte": fromDate,
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1, err2)
		}
	} else if (input.StartDate != nil) {
		fromDate, err1 := FormatDate(*input.StartDate, false, offset)

		if err1 == nil {
			filter["createdAt"] = map[string]interface{}{
				"$gte": fromDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1)
		}
	} else if (input.EndDate != nil) {
		toDate, err2 := FormatDate(*input.EndDate, true, offset)

		if err2 == nil {
			filter["createdAt"] = map[string]interface{}{
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err2)
		}
	}

	if input.ProcessStatus != nil {
		// Split the comma-separated string
		processStatusArray := strings.Split(*input.ProcessStatus, ",")
		// Trim spaces and collect unique statuses as integers
		var trimmedArray []int
		uniqueStatuses := make(map[int]bool)
		for _, status := range processStatusArray {
			trimmedStatus := strings.TrimSpace(status)
			if trimmedStatus != "" {
				statusInt, err := strconv.Atoi(trimmedStatus)
				if err == nil && !uniqueStatuses[statusInt] {
					trimmedArray = append(trimmedArray, statusInt)
					uniqueStatuses[statusInt] = true
				}
			}
		}
		if len(trimmedArray) > 0 {
			filter["processStatus"] = map[string]interface{}{
				"$in": trimmedArray,
			}
		}
	}

	if input.ForecastWeightFrom != nil && input.ForecastWeightTo != nil {
		filter["forecastWeight"] = map[string]interface{}{
			"$gte": *input.ForecastWeightFrom,
			"$lte": *input.ForecastWeightTo,
		}
	} else if (input.ForecastWeightFrom != nil) {
		filter["forecastWeight"] = map[string]interface{}{
			"$gte": *input.ForecastWeightFrom,
		}
	} else if (input.ForecastWeightTo != nil) {
		filter["forecastWeight"] = map[string]interface{}{
			"$gte": 0,
			"$lte": *input.ForecastWeightTo,
		}
	}
	
	filter["docType"] = "packing"

	filterJSON, err := json.Marshal(filter)
	if err != nil {
		fmt.Printf("Error marshalling filter to JSON: %v\n", err)
	} else {
		fmt.Printf("packing filter: %s\n", filterJSON)
	}

	return filter
}

func PackingFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllPacking, filter map[string]interface{}) ([]*models.PackingTransactionResponse, map[string]interface{}, error) {
	search, searchExists := filter["search"]

	if searchExists {
		delete(filter, "search")
	}

	// Initialize the base selector
	selector := map[string]interface{}{
		"selector": filter,
	}

	if searchExists && search != "" {
		selector["selector"] = map[string]interface{}{
			"$and": []map[string]interface{}{
				filter,
				{
					"$or": []map[string]interface{}{
						{"gmp": map[string]interface{}{"$regex": search}},
						{"packingHouseName": map[string]interface{}{"$regex": search}},
						{"gap": map[string]interface{}{"$regex": search}},
					},
				},
			},
		}
	} 

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getStringPacking, err := json.Marshal(selector)
	if err != nil {
		return nil, selector, err
	}

	fmt.Printf("Packing %s query\n", getStringPacking)

	queryPacking, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringPacking), int32(input.Limit), "")
	if err != nil {
		return nil, selector, err
	}
	defer queryPacking.Close()

	var dataPacking []*models.PackingTransactionResponse
	for queryPacking.HasNext() {
		queryResponse, err := queryPacking.Next()
		if err != nil {
			return nil, selector, err
		}

		var asset models.PackingTransactionResponse
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, selector, err
		}

		dataPacking = append(dataPacking, &asset)
	}

	return dataPacking, selector, nil
}