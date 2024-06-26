package chaincode

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreateGMP(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityGmp := models.TransactionGmp{}
	inputInterface, err := utils.Unmarshal(args, entityGmp)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionGmp)

	// err := ctx.GetClientIdentity().AssertAttributeValue("gmp.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have gmp.creator role1")
	}

	existsGmp, err := utils.AssetExists(ctx, input.Id)
	utils.HandleError(err)
	if existsGmp {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	TimeGmp := utils.GetTimeNow()

	asset := models.TransactionGmp{
		Id:                         input.Id,
		PackerId: 									input.PackerId,		
		PackingHouseRegisterNumber: input.PackingHouseRegisterNumber,
		Address:                    input.Address,
		PackingHouseName:           input.PackingHouseName,
		UpdatedDate:                input.UpdatedDate,
		Source:                     input.Source,
		Owner:                      clientID,
		OrgName:                    orgName,
		UpdatedAt:                  TimeGmp,
		CreatedAt:                  TimeGmp,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateGmp(ctx contractapi.TransactionContextInterface, args string) error {

	entityGmp := models.TransactionGmp{}
	inputInterface, err := utils.Unmarshal(args, entityGmp)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionGmp)

	asset, err := s.ReadGmp(ctx, input.Id)
	utils.HandleError(err)

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)
	if clientID != asset.Owner {
		return utils.ReturnError(utils.UNAUTHORIZE)
	}

	UpdatedGmp := utils.GetTimeNow()

	asset.Id = input.Id
	asset.PackerId = input.PackerId
	asset.PackingHouseRegisterNumber = input.PackingHouseRegisterNumber
	asset.Address = input.Address
	asset.PackingHouseName = input.PackingHouseName
	asset.UpdatedDate = input.UpdatedDate
	asset.Source = input.Source
	asset.UpdatedAt = UpdatedGmp

	assetJSON, errG := json.Marshal(asset)
	utils.HandleError(errG)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteGmp(ctx contractapi.TransactionContextInterface, id string) error {

	assetGmp, err := s.ReadGmp(ctx, id)
	utils.HandleError(err)

	clientIDGmp, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientIDGmp != assetGmp.Owner {
		return utils.ReturnError(utils.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) ReadGmp(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionGmp, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionGmp
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) GetAllGMP(ctx contractapi.TransactionContextInterface, args string) (*models.GmpGetAllResponse, error) {

	entityGetAllGmp := models.FilterGetAllGmp{}
	interfaceGmp, err := utils.Unmarshal(args, entityGetAllGmp)
	if err != nil {
		return nil, err
	}
	inputGmp := interfaceGmp.(*models.FilterGetAllGmp)
	filterGmp := utils.GmpSetFilter(inputGmp)

	queryStringGmp, err := utils.BuildQueryString(filterGmp)
	if err != nil {
		return nil, err
	}

	total, err := utils.CountTotalResults(ctx, queryStringGmp)
	if err != nil {
		return nil, err
	}

	if inputGmp.Skip > total {
		return nil, fmt.Errorf(utils.SKIPOVER)
	}

	assets, err := utils.GmpFetchResultsWithPagination(ctx, inputGmp, filterGmp)
	if err != nil {
		return nil, err
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].UpdatedAt.Before(assets[j].UpdatedAt)
	})

	if len(assets) == 0 {
		assets = []*models.GmpTransactionResponse{}
	}

	return &models.GmpGetAllResponse{
		Data:  "All Gmp",
		Obj:   assets,
		Total: total,
	}, nil
}

func (s *SmartContract) GetGmpByPackingHouseNumber(ctx contractapi.TransactionContextInterface, packingHouseRegisterNumber string) (*models.GetByRegisterNumberResponse, error) {
	// Get the asset using CertID
	queryKeyPackingHouse := fmt.Sprintf(`{"selector":{"packingHouseRegisterNumber":"%s"}}`, packingHouseRegisterNumber)

	resultsIteratorPackingHouse, err := ctx.GetStub().GetQueryResult(queryKeyPackingHouse)
	var asset *models.GmpTransactionResponse
	resData := "Get gmp by packingHouseRegisterNumber"
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIteratorPackingHouse.Close()

	if !resultsIteratorPackingHouse.HasNext() {
		resData = "Not found gmp by packingHouseRegisterNumber"

		return &models.GetByRegisterNumberResponse{
			Data: resData,
			Obj:  asset,
		}, nil
	}

	queryResponse, err := resultsIteratorPackingHouse.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &models.GetByRegisterNumberResponse{
		Data: resData,
		Obj:  asset,
	}, nil

}

func (s *SmartContract) FilterGmp(ctx contractapi.TransactionContextInterface, key, value string) ([]*models.TransactionGmp, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetGmp []*models.TransactionGmp
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.TransactionGmp
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetGmp = append(assetGmp, &asset)
		}
	}

	sort.Slice(assetGmp, func(i, j int) bool {
		return assetGmp[i].UpdatedAt.After(assetGmp[j].UpdatedAt)
	})

	return assetGmp, nil
}

func (s *SmartContract) CreateGmpCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.TransactionGmp

	errInputGmp := json.Unmarshal([]byte(args), &inputs)
	utils.HandleError(errInputGmp)

	for _, input := range inputs {
		// err := ctx.GetClientIdentity().AssertAttributeValue("gmp.creator", "true")

		orgNameG, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		existGmp, err := utils.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existGmp {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		clientIDG, err := utils.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		assetG := models.TransactionGmp{
			Id:                         input.Id,
			PackerId: 									input.PackerId,
			PackingHouseRegisterNumber: input.PackingHouseRegisterNumber,
			Address:                    input.Address,
			PackingHouseName:           input.PackingHouseName,
			UpdatedDate:                input.UpdatedDate,
			Source:                     input.Source,
			Owner:                      clientIDG,
			OrgName:                    orgNameG,
		}
		assetJSON, err := json.Marshal(assetG)
		if err != nil {
			return fmt.Errorf("failed to marshal asset JSON: %v", err)
		}

		err = ctx.GetStub().PutState(input.Id, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
		}

		fmt.Printf("Asset %s created successfully\n", input.Id)
	}

	return nil
}


func (s *SmartContract) UpdateMultipleGmp(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.TransactionGmp

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
		
		var existingAsset models.TransactionGmp
		err = json.Unmarshal(assetJSON, &existingAsset)
		if err != nil {
			return fmt.Errorf("failed to unmarshal existing asset: %v", err)
		}
		UpdatedGmp := utils.GetTimeNow()
		
		existingAsset.Id = input.Id
		existingAsset.PackerId = input.PackerId
		existingAsset.PackingHouseRegisterNumber = input.PackingHouseRegisterNumber
		existingAsset.Address = input.Address
		existingAsset.PackingHouseName = input.PackingHouseName
		existingAsset.UpdatedDate = input.UpdatedDate
		existingAsset.Source = input.Source
		existingAsset.UpdatedAt = UpdatedGmp

		updatedAssetJSON, err := json.Marshal(existingAsset)
		if err != nil {
			return fmt.Errorf("failed to marshal updated asset: %v", err)
		}
		
		err = ctx.GetStub().PutState(input.Id, updatedAssetJSON)
		if err != nil {
			return fmt.Errorf("failed to update asset in world state: %v", err)
		}
		
		fmt.Printf("Asset %s updated successfully\n", input.Id)
	}
	
	return nil
}