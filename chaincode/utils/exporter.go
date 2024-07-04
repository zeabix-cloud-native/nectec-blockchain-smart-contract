package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func ExporterSetFilter(input *models.ExporterFilterGetAll) map[string]interface{} {
    var filter = map[string]interface{}{}

    filter["docType"] = "exporter"

    if input.Province != nil {
        filter["province"] = *input.Province
    }

    if input.District != nil {
        filter["district"] = *input.District
    }

    if input.IssueDateFrom != nil && input.IssueDateTo != nil {
        filter["issueDate"] = map[string]interface{}{
            "$gte": *input.IssueDateFrom,
            "$lte": *input.IssueDateTo,
        }
    }

    if input.ExpireDateFrom != nil && input.ExpireDateTo != nil {
        filter["expiredDate"] = map[string]interface{}{
            "$gte": *input.ExpireDateFrom,
            "$lte": *input.ExpireDateTo,
        }
    }

    return filter
}

func ExporterFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.ExporterFilterGetAll, filter map[string]interface{}) ([]*models.ExporterTransactionResponse, int, error) {
	search := input.Search

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
                        {"plantType": map[string]interface{}{"$regex": searchTerm}},
                        {"name": map[string]interface{}{"$regex": searchTerm}},
                    },
                },
            },
        }
    }

    getStringE, err := json.Marshal(selector)
    if err != nil {
        return nil, 0, err
    }

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

