package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreatePlantTypeCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.PlantTypeModel 
	var eventPayloads []models.PlantTypeModel

	// Unmarshal input arguments
	errInputPlantType := json.Unmarshal([]byte(args), &inputs)
	if errInputPlantType != nil {
		return fmt.Errorf("failed to unmarshal input arguments: %v", errInputPlantType)
	}

	// Process each input
	for _, input := range inputs {
		// Get client's MSP ID
		orgNameG, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		// Check if the asset already exists
		existPlantType, err := utils.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existPlantType {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		// Get client's identity
		clientIDG, err := utils.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		// Create the asset
		assetG := models.PlantTypeModel{
			Id:          input.Id,
			Name:   	 input.Name,
			Address:   	 input.Address,
			Province:   	 input.Province,
			District:   	 input.District,
			PostCode:   	 input.PostCode,
			Email:   	 input.Email,
			IssueDate:   	 input.IssueDate,
			ExpiredDate:   	 input.ExpiredDate,
			PlantType:   input.PlantType,
			ExporterId:  input.ExporterId,
			Owner:       clientIDG,
			OrgName:     orgNameG,
			CreatedAt:   input.CreatedAt,
			DocType: 	 models.PlantType,
			UpdatedAt:   input.UpdatedAt,
		}

		// Marshal the asset to JSON
		assetJSON, err := json.Marshal(assetG)
		if err != nil {
			return fmt.Errorf("failed to marshal asset JSON: %v", err)
		}

		// Put the asset state
		err = ctx.GetStub().PutState(input.Id, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
		}

		// Add the asset to the event payloads
		eventPayloads = append(eventPayloads, assetG)
		fmt.Printf("Asset %s created successfully\n", input.Id)
	}

	return nil
}

func (s *SmartContract) ReadPlanType(ctx contractapi.TransactionContextInterface, id string) (*models.PlantTypeModel, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.PlantTypeModel
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	asset.IsCanDelete = true

	if (asset.ExporterId != "") {
		queryStr := fmt.Sprintf(`{
			"selector": {
				"docType": "formE",
				"createdById": "%s"
			},
			"use_index": [
				"_design/index-CreatedAt",
				"index-CreatedAt"
			]
		}`, asset.ExporterId)

		salesResultsIterator, err := ctx.GetStub().GetQueryResult(queryStr)
		if err != nil {
			return nil, fmt.Errorf("failed to query related sales: %v", err)
		}
		defer salesResultsIterator.Close()

		// If there are any related sales, set isCanDelete to false
		if salesResultsIterator.HasNext() {
			asset.IsCanDelete = false
		}
	}

	return &asset, nil
}

func (s *SmartContract) UpdatePlantType(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityPlanType := models.PlantTypeModel{}
	inputInterface, err := utils.Unmarshal(args, entityPlanType)
	utils.HandleError(err)
	input := inputInterface.(*models.PlantTypeModel)

	asset, err := s.ReadPlanType(ctx, input.Id)
	utils.HandleError(err)

	asset.UpdatedAt = input.UpdatedAt
	asset.ExporterId = input.ExporterId

	assetJSON, errE := json.Marshal(asset)
	utils.HandleError(errE)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) QueryPlanTypeWithPagination(ctx contractapi.TransactionContextInterface, filterParams string) (*models.PlantTypeResponse, error) {
	var filters models.PlanTypeFilterParams
    const offset = 7 // UTC+7

	err := json.Unmarshal([]byte(filterParams), &filters)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter parameters: %v", err)
	}

	selector := map[string]interface{}{
		"docType": "plantType",
	}

	if filters.AvailablePlanType != "" {
		selector["exporterId"] = ""
	}

	if filters.PlantType != "" {
		selector["plantType"] = filters.PlantType
	}

	if filters.Province != nil {
        selector["province"] = map[string]interface{}{
            "$regex": filters.Province,
        }
    }

    if filters.District != nil {
        selector["district"] = map[string]interface{}{
            "$regex": filters.District,
        }
    }

	if filters.Search != nil {
        selector["plantType"] = map[string]interface{}{
            "$regex": filters.Search,
        }
    }

	if filters.CreatedAtFrom != nil && filters.CreatedAtTo != nil {
		fromDate, err1 := utils.FormatDate(*filters.CreatedAtFrom, false, offset)
		toDate, err2 := utils.FormatDate(*filters.CreatedAtTo, true, offset)

		if err1 == nil && err2 == nil {
			selector["createdAt"] = map[string]interface{}{
				"$gte": fromDate,
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1, err2)
		}
	} else if (filters.CreatedAtFrom != nil) {
		fromDate, err1 := utils.FormatDate(*filters.CreatedAtFrom, false, offset)

		if err1 == nil {
			selector["createdAt"] = map[string]interface{}{
				"$gte": fromDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1)
		}
	} else if (filters.CreatedAtTo != nil) {
		toDate, err2 := utils.FormatDate(*filters.CreatedAtTo, true, offset)

		if err2 == nil {
			selector["createdAt"] = map[string]interface{}{
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err2)
		}
	}

	if filters.ExpireDateFrom != nil && filters.ExpireDateTo != nil {
		fromDate, err1 := utils.FormatDate(*filters.ExpireDateFrom, false, offset)
		toDate, err2 := utils.FormatDate(*filters.ExpireDateTo, true, offset)

		if err1 == nil && err2 == nil {
			selector["expiredDate"] = map[string]interface{}{
				"$gte": fromDate,
				"$lte": toDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1, err2)
		}
	} else if (filters.ExpireDateFrom != nil) {
		fromDate, err1 := utils.FormatDate(*filters.ExpireDateFrom, false, offset)

		if err1 == nil {
			selector["expiredDate"] = map[string]interface{}{
				"$gte": fromDate,
			}
		} else {
			fmt.Printf("Error formatting issue dates: %v, %v\n", err1)
		}
	} else if (filters.ExpireDateTo != nil) {
		toDate, err2 := utils.FormatDate(*filters.ExpireDateTo, true, offset)

		if err2 == nil {
			selector["expiredDate"] = map[string]interface{}{
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
		return &models.PlantTypeResponse{
			Data:  []*models.PlantTypeModel{},
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

	fmt.Printf("PlantType query %v", queryString)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal query string: %v", err)
	}

	// Execute the paginated query
	resultsIterator, _, err := ctx.GetStub().GetQueryResultWithPagination(string(queryString), int32(filters.Limit), "")
	if err != nil {
		return nil, fmt.Errorf("failed to get query result with pagination: %v", err)
	}
	defer resultsIterator.Close()

	var assets []*models.PlantTypeModel
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var asset models.PlantTypeModel
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		assets = append(assets, &asset)
	}

	for _, asset := range assets {
		asset.IsCanDelete = true

		if (asset.ExporterId != "") {
			queryStr := fmt.Sprintf(`{
				"selector": {
					"docType": "formE",
					"createdById": "%s"
				},
				"use_index": [
					"_design/index-CreatedAt",
					"index-CreatedAt"
				]
			}`, asset.ExporterId)
	
			salesResultsIterator, err := ctx.GetStub().GetQueryResult(queryStr)
			if err != nil {
				return nil, fmt.Errorf("failed to query related sales: %v", err)
			}
			defer salesResultsIterator.Close()
	
			// If there are any related sales, set isCanDelete to false
			if salesResultsIterator.HasNext() {
				asset.IsCanDelete = false
			}
		}
	}

	return &models.PlantTypeResponse{
		Data:  assets,
		Total: total,
	}, nil
}

func (s *SmartContract) GetPlantTypeByPlantType(ctx contractapi.TransactionContextInterface, plantType string) (*models.PlantTypeModel, error) {
	// Get the asset using CertID
	queryKeyPlantType := fmt.Sprintf(`{"selector":{"plantType":"%s", "docType":"plantType"}}`, plantType)

	resultsIteratorPlantType, err := ctx.GetStub().GetQueryResult(queryKeyPlantType)
	var asset *models.PlantTypeModel
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIteratorPlantType.Close()

	if !resultsIteratorPlantType.HasNext() {
		return &models.PlantTypeModel{}, nil
	}

	queryResponse, err := resultsIteratorPlantType.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return asset, nil
}

func (s *SmartContract) DeletePlantType(ctx contractapi.TransactionContextInterface, id string) error {
	plantType, err := s.ReadPlanType(ctx, id)
	utils.HandleError(err)

	return ctx.GetStub().DelState(plantType.Id)
}

func (s *SmartContract) UpdateMultiplePlantType(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.PlantTypeModel

	errInputGap := json.Unmarshal([]byte(args), &inputs)
	utils.HandleError(errInputGap)
	
	for _, input := range inputs {
		assetJSON, err := ctx.GetStub().GetState(input.Id)
		if err != nil {
			return fmt.Errorf("failed to read from world state: %v", err)
		}
		if assetJSON == nil {
			return fmt.Errorf("asset with ID %s does not exist", input.Id)
		}
		
		var existingAsset models.PlantTypeModel
		err = json.Unmarshal(assetJSON, &existingAsset)
		if err != nil {
			return fmt.Errorf("failed to unmarshal existing asset: %v", err)
		}

		existingAsset.Name =   input.Name
		existingAsset.Address =   input.Address
		existingAsset.Province =   input.Province
		existingAsset.District =    input.District
		existingAsset.PostCode =    input.PostCode
		existingAsset.Email =    input.Email
		existingAsset.IssueDate =   input.IssueDate
		existingAsset.ExpiredDate =   input.ExpiredDate
		existingAsset.PlantType =   input.PlantType
		existingAsset.Province =    input.Province
		existingAsset.UpdatedAt = 	input.UpdatedAt
		
		updatedAssetJSON, err := json.Marshal(existingAsset)
		if err != nil {
			return fmt.Errorf("failed to marshal updated asset: %v", err)
		}
		
		err = ctx.GetStub().PutState(input.Id, updatedAssetJSON)
		if err != nil {
			return fmt.Errorf("failed to update asset in world state: %v", err)
		}
		
		fmt.Printf("PlantType Asset %s updated successfully\n", input.Id)
	}
	
	return nil
}

func (s *SmartContract) GetPlantTypeList(ctx contractapi.TransactionContextInterface, args string) ([]models.PlantTypeModel, error) {
	var inputs []string

	errInputPlantType := json.Unmarshal([]byte(args), &inputs)
	if errInputPlantType != nil {
		return nil, fmt.Errorf("failed to unmarshal input arguments: %v", errInputPlantType)
	}

	var plantTypes []models.PlantTypeModel

	for _, plantType := range inputs {
		queryString := fmt.Sprintf(`{
			"selector": {
				"docType": "plantType",
				"plantType": "%s"
			}
		}`, plantType)
	
		resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
		if err != nil {
			return nil, fmt.Errorf("failed to execute rich query: %v", err)
		}
		defer resultsIterator.Close()
	
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, fmt.Errorf("failed to iterate results: %v", err)
			}
	
			var plantTypeModel models.PlantTypeModel
			err = json.Unmarshal(queryResponse.Value, &plantTypeModel)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal query response: %v", err)
			}
	
			plantTypes = append(plantTypes, plantTypeModel)
		}
	}

	// Ensure slices are initialized to empty slices if they are nil
	if plantTypes == nil {
		plantTypes = []models.PlantTypeModel{}
	}
	
	return plantTypes, nil
}