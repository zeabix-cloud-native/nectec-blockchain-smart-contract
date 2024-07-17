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
	countResultsIterator, err := ctx.GetStub().GetQueryResult(string(countQueryString))
	if err != nil {
		return nil, fmt.Errorf("failed to get count query result: %v", err)
	}
	defer countResultsIterator.Close()

	// Count the total number of records
	var totalCount int
	for countResultsIterator.HasNext() {
		_, err := countResultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next count query result: %v", err)
		}
		totalCount++
	}

	// If no records are found, return an empty response
	if totalCount == 0 {
		return &models.PlantTypeResponse{
			Data:  []*models.PlantTypeModel{},
			Total: totalCount,
		}, nil
	}

	// Create query string for paginated results
	queryString, err := json.Marshal(map[string]interface{}{
		"selector": selector,
		"sort": []map[string]string{
			{"createdAt": "desc"},
		},
        "use_index": []string{
            "_design/index-CreatedAt",
            "index-CreatedAt",
        },
	})

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

	return &models.PlantTypeResponse{
		Data:  assets,
		Total: totalCount,
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