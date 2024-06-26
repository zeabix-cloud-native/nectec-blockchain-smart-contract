package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreateFarmerProfile(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityFarmer := models.TransactionFarmer{}
	inputInterface, err := utils.Unmarshal(args, entityFarmer)

	utils.HandleError(err)
	input := inputInterface.(*models.TransactionFarmer)

	// err := ctx.GetClientIdentity().AssertAttributeValue("farmer.creator", "true")
	// if err != nil {
	// 	return fmt.Errorf("submitting client not authorized to create asset, does not have abac.creator role")
	// }

	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have farmer.creator role")
	}

	exists, err := utils.AssetExists(ctx, input.Id)
	utils.HandleError(err)
	if exists {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	CreatedAt := utils.GetTimeNow()

	asset := models.TransactionFarmer{
		Id:        input.Id,
		CertId:    input.CertId,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: CreatedAt,
		CreatedAt: CreatedAt,
		FarmerGaps: input.FarmerGaps,
		DocType: models.Farmer,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateFarmerProfile(ctx contractapi.TransactionContextInterface,
	args string) error {
	entityType := models.TransactionFarmer{}
	inputInterface, err := utils.Unmarshal(args, entityType)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionFarmer)

	asset, err := s.ReadFarmerProfile(ctx, input.Id)
	utils.HandleError(err)

	UpdatedAt := utils.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedAt
	asset.FarmerGaps = input.FarmerGaps

	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteFarmerProfile(ctx contractapi.TransactionContextInterface, id string) error {

	asset, err := s.ReadFarmerProfile(ctx, id)
	utils.HandleError(err)

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientID != asset.Owner {
		return fmt.Errorf(utils.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) ReadFarmerProfile(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionFarmer, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionFarmer
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	
	if asset.FarmerGaps == nil {
		asset.FarmerGaps = []models.FarmerGap{} 
	}

	return &asset, nil
}

func (s *SmartContract) GetAllFarmerProfile(ctx contractapi.TransactionContextInterface, args string) (*models.FarmerGetAllResponse, error) {

	var filter = map[string]interface{}{}
	filter["docType"] = "farmer"

	entityGetAll := models.FilterGetAllFarmer{}
	inputInterface, err := utils.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := inputInterface.(*models.FilterGetAllFarmer)

	queryString, err := utils.BuildQueryString(filter)
	if err != nil {
		return nil, err
	}

	total, err := utils.CountTotalResults(ctx, queryString)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, fmt.Errorf(utils.SKIPOVER)
	}

	arrFarmer, err := utils.FarmerFetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrFarmer, func(i, j int) bool {
		return arrFarmer[i].UpdatedAt.Before(arrFarmer[j].UpdatedAt)
	})

	if len(arrFarmer) == 0 {
		arrFarmer = []*models.FarmerTransactionResponse{}
	}
	
	for _, farmer := range arrFarmer {
		log.Printf("farmer item %v", farmer)
	}

	return &models.FarmerGetAllResponse{
		Data:  "All Farmer",
		Obj:   arrFarmer,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterFarmer(ctx contractapi.TransactionContextInterface, key, value string) ([]*models.TransactionFarmer, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetFarmer []*models.TransactionFarmer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.TransactionFarmer
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetFarmer = append(assetFarmer, &asset)
		}
	}

	sort.Slice(assetFarmer, func(i, j int) bool {
		return assetFarmer[i].UpdatedAt.After(assetFarmer[j].UpdatedAt)
	})
	return assetFarmer, nil
}

func (s *SmartContract) GetFarmerHistory(ctx contractapi.TransactionContextInterface, key string) ([]*models.FarmerTransactionHistory, error) {
	// Get the history for the specified key
	farmerResultsIterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for key %s: %v", key, err)
	}
	defer farmerResultsIterator.Close()

	var farmerHistory []*models.FarmerTransactionHistory
	var farmerAssets []*models.FarmerTransactionResponse
	for farmerResultsIterator.HasNext() {
		// Get the next history record
		farmer, err := farmerResultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next history farmer for key %s: %v", key, err)
		}

		var asset models.FarmerTransactionResponse
		if !farmer.IsDelete {
			err = json.Unmarshal(farmer.Value, &asset)
			if err != nil {
				return nil, err
			}
			farmerAssets = append(farmerAssets, &asset)

		} else {
			farmerAssets = []*models.FarmerTransactionResponse{}
		}
		// Convert the timestamp to string in the desired format
		timestampStr := time.Unix(farmer.Timestamp.Seconds, int64(farmer.Timestamp.Nanos)).Format(utils.TIMEFORMAT)

		historyFarmer := &models.FarmerTransactionHistory{
			TxId:      farmer.TxId,
			Value:     farmerAssets,
			Timestamp: timestampStr,
			IsDelete:  farmer.IsDelete,
		}

		farmerHistory = append(farmerHistory, historyFarmer)
	}

	return farmerHistory, nil
}

func (s *SmartContract) GetLastIdFarmer(ctx contractapi.TransactionContextInterface) string {
	// Query to get all records sorted by ID in descending order
	query := `{
			"selector": {
				"docType": "farmer"
			},
			"sort": [{"_id": "desc"}],
			"limit": 1,
			"use_index": "index-id"
	}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return "error querying CouchDB"
	}
	defer resultsIterator.Close()

	// Check if there is a result
	if !resultsIterator.HasNext() {
		return ""
	}

	// Get the first (and only) result
	queryResponse, err := resultsIterator.Next()
	if err != nil {
		return "error iterating query results"
	}

	var result struct {
		Id string `json:"id"`
	}

	// Unmarshal the result into the result struct
	if err := json.Unmarshal(queryResponse.Value, &result); err != nil {
		return "error unmarshalling document"
	}

	return result.Id
}

func (s *SmartContract) CreateFarmerFromCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.TransactionFarmer
	var eventPayloads []models.TransactionFarmer

	errInput := json.Unmarshal([]byte(args), &inputs)
	if errInput != nil {
		return fmt.Errorf("failed to unmarshal JSON array: %v", errInput)
	}

	for _, input := range inputs {
		orgName, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		existFarmer, err := utils.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existFarmer {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		clientID, err := utils.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		asset := models.TransactionFarmer{
			Id:        input.Id,
			CertId:    input.CertId,
			FarmerGaps: input.FarmerGaps,
			Owner:     clientID,
			OrgName:   orgName,
			UpdatedAt: input.CreatedAt,
			CreatedAt: input.UpdatedAt,
			DocType: models.Farmer,
		}

		assetJSON, err := json.Marshal(asset)
		eventPayloads = append(eventPayloads, asset)
		if err != nil {
			return fmt.Errorf("failed to marshal asset JSON: %v", err)
		}

		err = ctx.GetStub().PutState(input.Id, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
		}

		fmt.Printf("Asset %s created successfully\n", input.Id)

	}

	eventPayloadJSON, err := json.Marshal(eventPayloads)
	if err != nil {
		return fmt.Errorf("failed to marshal asset JSON: %v", err)
	}
	ctx.GetStub().SetEvent("batchCreatedUserEvent", eventPayloadJSON)

	return nil
}