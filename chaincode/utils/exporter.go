package utils

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func ExporterSetFilter(input *models.ExporterFilterGetAll) map[string]interface{} {
    var filter = map[string]interface{}{}
    const offset = 7 // UTC+7
    
    filter["docType"] = "exporter"

    if input.Province != nil {
        filter["plantTypeDetail.province"] = map[string]interface{}{
            "$regex": input.Province,
        }
    }

    if input.District != nil {
        filter["plantTypeDetail.district"] = map[string]interface{}{
            "$regex": input.District,
        }
    }

    if input.CreatedAtFrom != nil && input.CreatedAtTo != nil {
		fromDate, err1 := FormatDate(*input.CreatedAtFrom, false, offset)
		toDate, err2 := FormatDate(*input.CreatedAtTo, true, offset)

		if err1 == nil && err2 == nil {
			filter["plantTypeDetail.createdAt"] = map[string]interface{}{
				"$gte": fromDate,
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1, err2)
		}
	} else if (input.CreatedAtFrom != nil) {
		fromDate, err1 := FormatDate(*input.CreatedAtFrom, false, offset)

		if err1 == nil {
			filter["plantTypeDetail.createdAt"] = map[string]interface{}{
				"$gte": fromDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1)
		}
	} else if (input.CreatedAtTo != nil) {
		toDate, err2 := FormatDate(*input.CreatedAtTo, true, offset)

		if err2 == nil {
			filter["plantTypeDetail.createdAt"] = map[string]interface{}{
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
			filter["plantTypeDetail.expiredDate"] = map[string]interface{}{
				"$gte": fromDate,
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1, err2)
		}
	} else if (input.ExpireDateFrom != nil) {
		fromDate, err1 := FormatDate(*input.ExpireDateFrom, false, offset)

		if err1 == nil {
			filter["plantTypeDetail.expiredDate"] = map[string]interface{}{
				"$gte": fromDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1)
		}
	} else if (input.ExpireDateTo != nil) {
		toDate, err2 := FormatDate(*input.ExpireDateTo, true, offset)

		if err2 == nil {
			filter["plantTypeDetail.expiredDate"] = map[string]interface{}{
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err2)
		}
	}

    return filter
}

func ExporterFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.ExporterFilterGetAll, filter map[string]interface{}) ([]*models.ExporterTransactionResponse, int, error) {
	search := input.Search

    selector := map[string]interface{}{
        "selector": filter,
        "use_index": []string{
			"_design/index-CreatedAt",
			"index-CreatedAt",
		},
    }

    if search != nil && *search != "" {
        searchTerm := *search
        selector["selector"] = map[string]interface{}{
            "$and": []map[string]interface{}{
                filter,
                {
                    "$or": []map[string]interface{}{
                        {"id": map[string]interface{}{"$regex": searchTerm}},
                        {"plantTypeDetail.plantType": map[string]interface{}{"$regex": searchTerm}},
                        {"plantTypeDetail.name": map[string]interface{}{"$regex": searchTerm}},
                    },
                },
            },
        }
    }

    getStringE, err := json.Marshal(selector)
    if err != nil {
        return nil, 0, err
    }

    fmt.Printf("Query String for Fetching Results: %s\n", getStringE) // Debugging

    // Fetch total count of the results
    resultsIterator, err := ctx.GetStub().GetQueryResult(string(getStringE))
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

    // Apply pagination
    if input.Skip > 0 {
        selector["skip"] = input.Skip
    }
    if input.Limit > 0 {
        selector["limit"] = input.Limit
    }

    getStringE, err = json.Marshal(selector)
    if err != nil {
        return nil, 0, err
    }
	fmt.Printf("queryString %s", getStringE)

    queryGmp, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringE), int32(input.Limit), "")
    if err != nil {
        return nil, 0, err
    }
    defer queryGmp.Close()

    var dataGmp []*models.ExporterTransactionResponse
    for queryGmp.HasNext() {
        queryRes, err := queryGmp.Next()
        if err != nil {
            return nil, 0, err
        }

        var dataG models.ExporterTransactionResponse
        err = json.Unmarshal(queryRes.Value, &dataG)
        if err != nil {
            return nil, 0, err
        }

        dataGmp = append(dataGmp, &dataG)
    }

    return dataGmp, total, nil
}

