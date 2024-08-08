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
		filter["farmerId"] = input.FarmerID
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
	} else if (input.AreaRaiFrom != nil) {
		filter["areaRai"] = map[string]interface{}{
			"$gte": *input.AreaRaiFrom,
		}
	} else if (input.AreaRaiTo != nil) {
		filter["areaRai"] = map[string]interface{}{
			"$gte": 0,
			"$lte": *input.AreaRaiTo,
		}
	}

	if input.CreatedAtFrom != nil && input.CreatedAtTo != nil {
		fromDate, err1 := FormatDate(*input.CreatedAtFrom, false, offset)
		toDate, err2 := FormatDate(*input.CreatedAtTo, true, offset)

		if err1 == nil && err2 == nil {
			filter["createdAt"] = map[string]interface{}{
				"$gte": fromDate,
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1, err2)
		}
	} else if (input.CreatedAtFrom != nil) {
		fromDate, err1 := FormatDate(*input.CreatedAtFrom, false, offset)

		if err1 == nil {
			filter["createdAt"] = map[string]interface{}{
				"$gte": fromDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1)
		}
	} else if (input.CreatedAtTo != nil) {
		toDate, err2 := FormatDate(*input.CreatedAtTo, true, offset)

		if err2 == nil {
			filter["createdAt"] = map[string]interface{}{
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err2)
		}
	}

	if input.ExpireDateFrom != nil && input.ExpireDateTo != nil {
		fromDate, err1 := FormatDate(*input.ExpireDateFrom, false, offset)
		toDate, err2 := FormatDate(*input.ExpireDateTo, true, offset)

		if err1 == nil && err2 == nil {
			filter["expireDate"] = map[string]interface{}{
				"$gte": fromDate,
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1, err2)
		}
	} else if (input.ExpireDateFrom != nil) {
		fromDate, err1 := FormatDate(*input.ExpireDateFrom, false, offset)

		if err1 == nil {
			filter["expireDate"] = map[string]interface{}{
				"$gte": fromDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1)
		}
	} else if (input.ExpireDateTo != nil) {
		toDate, err2 := FormatDate(*input.ExpireDateTo, true, offset)

		if err2 == nil {
			filter["expireDate"] = map[string]interface{}{
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err2)
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

	if input.CertID != nil && *input.CertID != "" {
		searchTerm := *input.CertID
		selector["selector"] = map[string]interface{}{
			"$and": []map[string]interface{}{
				filter,
				{
					"$or": []map[string]interface{}{
						{"certId": map[string]interface{}{"$regex": searchTerm}},
						{"displayCertId": map[string]interface{}{"$regex": searchTerm}},
					},
				},
			},
		}
	}

	if input.Skip > 0 {
		selector["skip"] = input.Skip
	}
	if input.Limit > 0 {
		selector["limit"] = input.Limit
	}

	selector["sort"] = []map[string]string {
    	{"createdAt": "desc"},
	}

	queryString, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Query String for Fetching Results: %s\n", queryString) // Debugging

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


