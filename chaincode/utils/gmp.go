package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func GmpSetFilter(input *models.FilterGetAllGmp) map[string]interface{} {
	var filter = map[string]interface{}{}

	filter["docType"] = "gmp"

	return filter
}

func GmpFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.FilterGetAllGmp, filter map[string]interface{}) ([]*models.GmpTransactionResponse, int, error) {
    search := input.Search

    filter["docType"] = "gmp"

    selector := map[string]interface{}{
        "selector": filter,
    }

    if search != nil && *search != "" {
        searchTerm := *search
        selector["selector"] = map[string]interface{}{
            "$and": []map[string]interface{}{
                filter,
                {
                    "$or": []map[string]interface{}{
                        {"packingHouseRegisterNumber": map[string]interface{}{"$regex": searchTerm}},
                        {"packingHouseName": map[string]interface{}{"$regex": searchTerm}},
                        {"address": map[string]interface{}{"$regex": searchTerm}},
                    },
                },
            },
        }
    }

    getStringGmp, err := json.Marshal(map[string]interface{}{
		"selector": selector,
		"sort": []map[string]string{
			{"createdAt": "desc"},
		},
        "use_index": []string{
            "_design/index-CreatedAt",
            "index-CreatedAt",
        },
	})
    if err != nil {
        return nil, 0, err
    }

    // Fetch total count of the results
    total, err := getTotalGmpCount(ctx, string(getStringGmp))
    if err != nil {
        return nil, 0, err
    }

    // Apply pagination
    if input.Skip > 0 {
        selector["skip"] = input.Skip
    }
    if input.Limit > 0 {
        selector["limit"] = input.Limit
    }

    getStringGmpWithPagination, err := json.Marshal(selector)
    if err != nil {
        return nil, 0, err
    }

    queryGmp, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringGmpWithPagination), int32(input.Limit), "")
    if err != nil {
        return nil, 0, err
    }
    defer queryGmp.Close()

    var dataGmp []*models.GmpTransactionResponse
    for queryGmp.HasNext() {
        queryRes, err := queryGmp.Next()
        if err != nil {
            return nil, 0, err
        }

        var dataG models.GmpTransactionResponse
        err = json.Unmarshal(queryRes.Value, &dataG)
        if err != nil {
            return nil, 0, err
        }

        dataGmp = append(dataGmp, &dataG)
    }

    return dataGmp, total, nil
}

func getTotalGmpCount(ctx contractapi.TransactionContextInterface, query string) (int, error) {
    resultsIterator, err := ctx.GetStub().GetQueryResult(query)
    if err != nil {
        return 0, err
    }
    defer resultsIterator.Close()

    total := 0
    for resultsIterator.HasNext() {
        _, err := resultsIterator.Next()
        if err != nil {
            return 0, err
        }
        total++
    }
    return total, nil
}