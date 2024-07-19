package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s SmartContract) CreateRegulatorProfile(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {

	entityRegulator := models.TransactionRegulator{}
	inputInterface, err := utils.Unmarshal(args, entityRegulator)
	if err != nil {
		return err
	}

	input := inputInterface.(*models.TransactionRegulator)
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return utils.ReturnError(utils.UNAUTHORIZE)
	}

	existRegulator, err := utils.AssetExists(ctx, input.Id)
	utils.HandleError(err)
	if existRegulator {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	asset := models.TransactionRegulator{
		Id:        input.Id,
		CertId:    input.CertId,
		UserId:    input.UserId,
		Owner:     clientID,
		OrgName:   orgName,
		DocType:   models.Regulator,
		ProfileImg: input.ProfileImg,
		UpdatedAt: input.CreatedAt,
		CreatedAt: input.UpdatedAt,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s SmartContract) UpdateRegulatorProfile(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityRegulator := models.TransactionRegulator{}
	inputInterface, err := utils.Unmarshal(args, entityRegulator)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionRegulator)

	asset, err := s.ReadRegulatorProfile(ctx, input.Id)
	utils.HandleError(err)

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientID != asset.Owner {
		return utils.ReturnError(utils.UNAUTHORIZE)
	}

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = input.UpdatedAt
	asset.ProfileImg = input.ProfileImg
	asset.UserId = input.UserId

	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) ReadRegulatorProfile(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionRegulator, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionRegulator
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) QueryRegulatorWithPagination(ctx contractapi.TransactionContextInterface, filterParams string) (*models.RegulatorGetAllResponse, error) {
	var filters models.FilterGetAllRegulator
	err := json.Unmarshal([]byte(filterParams), &filters)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter parameters: %v", err)
	}

	selector := map[string]interface{}{
		"docType": "regulator",
	}

	if filters.UserId != "" {
		selector["userId"] = filters.UserId
	}
	
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
		return &models.RegulatorGetAllResponse{
			Obj:  []*models.TransactionRegulator{},
			Total: totalCount,
		}, nil
	}

	// Apply pagination
	if filters.Skip > 0 {
		selector["skip"] = filters.Skip
	}
	if filters.Limit > 0 {
		selector["limit"] = filters.Limit
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

	var assets []*models.TransactionRegulator
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var asset models.TransactionRegulator
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		assets = append(assets, &asset)
	}

	return &models.RegulatorGetAllResponse{
		Obj:  assets,
		Total: totalCount,
	}, nil
}

func (s *SmartContract) DeleteRegulator(ctx contractapi.TransactionContextInterface, id string) error {
	assetRegulator, err := s.ReadRegulatorProfile(ctx, id)
	utils.HandleError(err)

	// clientIDGmp, err := utils.GetIdentity(ctx)
	// utils.HandleError(err)

	// if clientIDGmp != assetRegulator.Owner {
	// 	return utils.ReturnError(utils.UNAUTHORIZE)
	// }

	return ctx.GetStub().DelState(assetRegulator.Id)
}

func (s *SmartContract) GetRegulatorByUserId(ctx contractapi.TransactionContextInterface, regulatorId string) (*models.RegulatorByIdResponse, error) {
	queryKeyFarmer := fmt.Sprintf(`{"selector":{"userId":"%s", "docType": "regulator"}}`, regulatorId)

	resultsIteratorFarmer, err := ctx.GetStub().GetQueryResult(queryKeyFarmer)
	var asset *models.TransactionRegulator
	resData := "Get regulator by regulatorId"
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIteratorFarmer.Close()

	if !resultsIteratorFarmer.HasNext() {
		resData = "Not found regulator by regulatorId"

		return &models.RegulatorByIdResponse{
			Data: resData,
			Obj:  asset,
		}, nil
	}

	queryResponse, err := resultsIteratorFarmer.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &models.RegulatorByIdResponse{
		Data: resData,
		Obj:  asset,
	}, nil
}