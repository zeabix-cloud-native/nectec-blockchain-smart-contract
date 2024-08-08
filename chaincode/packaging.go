package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreatePackagingCsv(
    ctx contractapi.TransactionContextInterface,
    args string,
) error {
    var inputs []models.TransactionPackaging

    errInputPackaging := json.Unmarshal([]byte(args), &inputs)
    if errInputPackaging != nil {
        return fmt.Errorf("error unmarshaling input: %v", errInputPackaging)
    }

    var batchErrors []error

    for _, input := range inputs {
        if err := s.processSinglePackaging(ctx, input); err != nil {
            batchErrors = append(batchErrors, err)
            fmt.Printf("Error processing asset %s: %v\n", input.Id, err)
        } else {
            fmt.Printf("Packaging Asset %s created successfully\n", input.Id)
        }
    }

    if len(batchErrors) > 0 {
        return fmt.Errorf("encountered errors during processing: %v", batchErrors)
    }

    return nil
}

func (s *SmartContract) processSinglePackaging(ctx contractapi.TransactionContextInterface, input models.TransactionPackaging) error {
    orgNamePackaging, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
    }

    existPackaging, err := utils.AssetExists(ctx, input.Id)
    if err != nil {
        return fmt.Errorf("error checking if asset exists: %v", err)
    }
    if existPackaging {
        return fmt.Errorf("the asset %s already exists", input.Id)
    }

    clientIDPackaging, err := utils.GetIdentity(ctx)
    if err != nil {
        return fmt.Errorf("failed to get submitting client's identity: %v", err)
    }

    // timestamp := utils.GenerateTimestamp()

    assetPackaging := models.TransactionPackaging{
        Id:          input.Id,
        ContainerId: input.ContainerId,
        ExportId:    input.ExportId,
        LotNumber:   input.LotNumber,
        BoxId:       input.BoxId,
        Gap:         input.Gap,
        Gmp:         input.Gmp,
        GradeName:   input.GradeName,
        Gtin14:      input.Gtin14,
        Gtin13:      input.Gtin13,
        VarietyName: input.VarietyName,
		CreatedById: input.CreatedById,
        DocType:     models.Packaging,
        Owner:       clientIDPackaging,
        OrgName:     orgNamePackaging,
        CreatedAt:   input.CreatedAt,
        UpdatedAt:   input.UpdatedAt,
    }

    assetJSON, err := json.Marshal(assetPackaging)
    if err != nil {
        return fmt.Errorf("failed to marshal asset JSON for asset %s: %v", input.Id, err)
    }

    err = ctx.GetStub().PutState(input.Id, assetJSON)
    if err != nil {
        return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
    }

    return nil
}

func (s *SmartContract) ReadPackaging(ctx contractapi.TransactionContextInterface, id string) (*models.ReadPackaging, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.ReadPackaging
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	gmp, err := s.GetPackagingGmp(ctx, asset.Gmp) 
	if err != nil {
		return nil, fmt.Errorf("failed to get gmp for asset %s: %v", id, err)
	}
	asset.GmpDetail = gmp

	return &asset, nil
}

func (s *SmartContract) GetPackagingGmp(ctx contractapi.TransactionContextInterface, id string) (models.TransactionGmp, error) {
    selector := map[string]interface{}{
        "docType": "gmp",
        "packingHouseRegisterNumber": id,
    }

    queryString, err := json.Marshal(map[string]interface{}{
        "selector": selector,
        "use_index": []string{
            "_design/index-docType",
            "index-docType",
        },
    })
    if err != nil {
        return models.TransactionGmp{}, fmt.Errorf("failed to marshal query string: %v", err)
    }

    fmt.Printf("GMP query: %s\n", queryString)

    resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryString))
    if err != nil {
        return models.TransactionGmp{}, fmt.Errorf("failed to get GMP query result: %v", err)
    }
    defer resultsIterator.Close()

    if !resultsIterator.HasNext() {
        return models.TransactionGmp{}, nil
    }

    queryResponse, err := resultsIterator.Next()
    if err != nil {
        return models.TransactionGmp{}, fmt.Errorf("failed to get next GMP query result: %v", err)
    }

    var gmpDocument models.TransactionGmp
    err = json.Unmarshal(queryResponse.Value, &gmpDocument)
    if err != nil {
        return models.TransactionGmp{}, fmt.Errorf("failed to unmarshal GMP query result: %v", err)
    }

    return gmpDocument, nil
}

func (s *SmartContract) QueryPackagingWithPagination(ctx contractapi.TransactionContextInterface, filterParams string) (*models.TransactionPackagingResponse, error) {
    const offset = 7 // UTC+7
	var filters models.PackagingFilterParams
	err := json.Unmarshal([]byte(filterParams), &filters)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter parameters: %v", err)
	}

	selector := map[string]interface{}{
		"docType": "packaging",
	}

	if filters.CreatedById != "" {
		selector["createdById"] = filters.CreatedById
	}

	if filters.ContainerId != "" {
		selector["containerId"] = filters.ContainerId
	}
	if filters.BoxId != "" {
		selector["boxId"] = filters.BoxId
	}
	if filters.Gtin13 != "" {
		selector["gtin13"] = filters.Gtin13
	}
	if filters.Gap != "" {
		selector["gap"] = filters.Gap
	}

	if filters.StartDate != nil && filters.EndDate != nil {
		fromDate, err1 := utils.FormatDate(*filters.StartDate, false, offset)
		toDate, err2 := utils.FormatDate(*filters.EndDate, true, offset)

		if err1 == nil && err2 == nil {
			selector["createdAt"] = map[string]interface{}{
				"$gte": fromDate,
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1, err2)
		}
	} else if (filters.StartDate != nil) {
		fromDate, err1 := utils.FormatDate(*filters.StartDate, false, offset)

		if err1 == nil {
			selector["createdAt"] = map[string]interface{}{
				"$gte": fromDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1)
		}
	} else if (filters.EndDate != nil) {
		toDate, err2 := utils.FormatDate(*filters.EndDate, true, offset)

		if err2 == nil {
			selector["createdAt"] = map[string]interface{}{
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err2)
		}
	}

	// Create query string for counting total records
	countQueryString, err := json.Marshal(map[string]interface{}{
		"selector": selector,
		"use_index": []string{
            "_design/index-CreatedAt",
            "index-CreatedAt",
        },
	})

	if err != nil {
		return nil, fmt.Errorf("failed to marshal count query string: %v", err)
	}

	// Execute the count query
	total, err := utils.CountTotalResults(ctx, string(countQueryString))
	if err != nil {
		return nil, err
	}

	// If no records are found, return an empty response
	if total == 0 {
		return &models.TransactionPackagingResponse{
			Data:  []*models.TransactionPackaging{},
			Total: total,
		}, nil
	}

	query := map[string]interface{}{
		"selector": selector,
		"sort": []map[string]string{
			{"createdAt": "desc"},
		},
        "use_index": []string{
            "_design/index-CreatedAt",
            "index-CreatedAt",
        },
	}

	if filters.Skip > 0 {
		query["skip"] = filters.Skip
	}
	if filters.Limit > 0 {
		query["limit"] = filters.Limit
	}

	// Create query string for paginated results
	queryString, err := json.Marshal(query)

	fmt.Printf("Packaging query %v", queryString)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal query string: %v", err)
	}

	// Execute the paginated query
	resultsIterator, _, err := ctx.GetStub().GetQueryResultWithPagination(string(queryString), int32(filters.Limit), "")
	if err != nil {
		return nil, fmt.Errorf("failed to get query result with pagination: %v", err)
	}
	defer resultsIterator.Close()

	var assets []*models.TransactionPackaging
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var asset models.TransactionPackaging
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		assets = append(assets, &asset)
	}

	return &models.TransactionPackagingResponse{
		Data:  assets,
		Total: total,
	}, nil
}


