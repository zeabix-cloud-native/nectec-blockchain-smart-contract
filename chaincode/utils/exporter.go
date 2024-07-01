package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
)

func ExporterSetFilter(input *models.ExporterFilterGetAll) map[string]interface{} {
	var filter = map[string]interface{}{}

	filter["docType"] = "exporter"

	if input.Search != nil {
		filter["search"] = *input.Search
	}

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

func ExporterFetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *models.ExporterFilterGetAll, filter map[string]interface{}) ([]*models.ExporterTransactionResponse, error) {
	search, searchExists := filter["search"]

	filter["docType"] = "exporter"

	if searchExists {
		delete(filter, "search")
	}

	selector := map[string]interface{}{
		"selector": filter,
	}

	if searchExists && search != "" {
		selector["selector"] = map[string]interface{}{
			"$and": []map[string]interface{}{
				filter,
				{
					"$or": []map[string]interface{}{
						{"plantType": map[string]interface{}{"$regex": search}},
						{"name": map[string]interface{}{"$regex": search}},
					},
				},
			},
		}
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getStringE, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryExporter, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringE), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryExporter.Close()

	var dataExporter []*models.ExporterTransactionResponse
	for queryExporter.HasNext() {
		queryRes, err := queryExporter.Next()
		if err != nil {
			return nil, err
		}

		var dataE models.ExporterTransactionResponse
		err = json.Unmarshal(queryRes.Value, &dataE)
		if err != nil {
			return nil, err
		}

		dataExporter = append(dataExporter, &dataE)
	}

	return dataExporter, nil
}
