package chaincode

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreatePacking(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityPacking := models.TransactionPacking{}
	inputInterface, err := utils.Unmarshal(args, entityPacking)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionPacking)

	// err := ctx.GetClientIdentity().AssertAttributeValue("packing.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have packing.creator role")
	}

	existsPacking, err := utils.AssetExists(ctx, input.Id)
	utils.HandleError(err)
	if existsPacking {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientIDPacking, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	TimePacking := utils.GetTimeNow()
	fmt.Printf("CreatedAt %v", TimePacking)

	asset := models.TransactionPacking{
		Id:             input.Id,
		OrderID:        input.OrderID,
		FarmerID:       input.FarmerID,
		PackingHouseName: input.PackingHouseName,
		ForecastWeight: input.ForecastWeight,
		ActualWeight:   input.ActualWeight,
		SavedTime:      input.SavedTime,
		ApprovedDate:   input.ApprovedDate,
		ApprovedType:   input.ApprovedType,
		FinalWeight:    input.FinalWeight,
		Remark:         input.Remark,
		PackerId:       input.PackerId,
		Gmp:            input.Gmp,
		Gap:            input.Gap,
		ProcessStatus:  input.ProcessStatus,
		SellingStep:  	input.SellingStep,
		Owner:          clientIDPacking,
		OrgName:        orgName,
		UpdatedAt:      TimePacking,
		CreatedAt:      TimePacking,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdatePacking(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityPacking := models.TransactionPacking{}
	inputInterface, err := utils.Unmarshal(args, entityPacking)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionPacking)

	asset, err := s.ReadPacking(ctx, input.Id)
	utils.HandleError(err)

	UpdatedPacking := utils.GetTimeNow()

	asset.Id = input.Id
	asset.OrderID = input.OrderID
	asset.FarmerID = input.FarmerID // not update
	asset.ForecastWeight = input.ForecastWeight
	asset.ActualWeight = input.ActualWeight
	asset.SavedTime = input.SavedTime
	asset.ApprovedDate = input.ApprovedDate
	asset.ApprovedType = input.ApprovedType
	asset.FinalWeight = input.FinalWeight
	asset.Remark = input.Remark
	asset.PackerId = input.PackerId // not update
	asset.Gmp = input.Gmp
	asset.Gap = input.Gap
	asset.ProcessStatus = input.ProcessStatus
	asset.SellingStep = input.SellingStep
	asset.UpdatedAt = UpdatedPacking

	assetJSON, errPacking := json.Marshal(asset)
	utils.HandleError(errPacking)

	ctx.GetStub().SetEvent("UpdateAsset", assetJSON)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeletePacking(ctx contractapi.TransactionContextInterface, id string) error {

	assetPacking, err := s.ReadPacking(ctx, id)
	if err != nil {
		return err
	}

	clientIDPacking, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientIDPacking != assetPacking.Owner {
		return utils.ReturnError(utils.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferPacking(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetPacking, err := s.ReadPacking(ctx, id)
	utils.HandleError(err)

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientID != assetPacking.Owner {
		return utils.ReturnError(utils.UNAUTHORIZE)
	}

	assetPacking.Owner = newOwner
	assetJSON, err := json.Marshal(assetPacking)
	utils.HandleError(err)
	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadPacking(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionPacking, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionPacking
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}
func (s *SmartContract) GetAllPacking(ctx contractapi.TransactionContextInterface, args string) (*models.PackingGetAllResponse, error) {

	entityGetAllPacking := models.FilterGetAllPacking{}
	interfacePacking, err := utils.Unmarshal(args, entityGetAllPacking)
	if err != nil {
		return nil, err
	}
	inputPacking := interfacePacking.(*models.FilterGetAllPacking)
	filterPacking := utils.PackingSetFilter(inputPacking)

	queryStringPacking, err := utils.BuildQueryString(filterPacking)
	if err != nil {
		return nil, err
	}

	total, err := utils.CountTotalResults(ctx, queryStringPacking)
	if err != nil {
		return nil, err
	}

	if inputPacking.Skip > total {
		return nil, fmt.Errorf(utils.SKIPOVER)
	}

	arrPacking, err := utils.PackingFetchResultsWithPagination(ctx, inputPacking, filterPacking)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrPacking, func(i, j int) bool {
		return arrPacking[i].UpdatedAt.Before(arrPacking[j].UpdatedAt)
	})

	if len(arrPacking) == 0 {
		arrPacking = []*models.PackingTransactionResponse{}
	}

	return &models.PackingGetAllResponse{
		Data:  "All Packing",
		Obj:   arrPacking,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterPacking(ctx contractapi.TransactionContextInterface, key, value string) ([]*models.TransactionPacking, error) {
	resultsIteratorP, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIteratorP.Close()

	var assetPacking []*models.TransactionPacking
	for resultsIteratorP.HasNext() {
		queryResponse, err := resultsIteratorP.Next()
		if err != nil {
			return nil, err
		}

		var asset models.TransactionPacking
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetPacking = append(assetPacking, &asset)
		}
	}

	sort.Slice(assetPacking, func(i, j int) bool {
		return assetPacking[i].UpdatedAt.After(assetPacking[j].UpdatedAt)
	})

	return assetPacking, nil
}

func (s *SmartContract) GetLatestHistoryForKey(ctx contractapi.TransactionContextInterface, key string) (*models.PackingTransactionHistory, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for key %s: %v", key, err)
	}
	defer resultsIterator.Close()

	var latestHistory *models.PackingTransactionHistory
	var packingAsset []*models.PackingTransactionResponse

	for resultsIterator.HasNext() {
		// Get the next history record
		record, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next history record for key %s: %v", key, err)
		}

		var asset models.PackingTransactionResponse
		if !record.IsDelete {
			err = json.Unmarshal(record.Value, &asset)
			if err != nil {
				return nil, err
			}
			packingAsset = append(packingAsset, &asset)
		} else {
			packingAsset = []*models.PackingTransactionResponse{}
		}

		// Sort packingAsset by SellingStep
		sort.SliceStable(packingAsset, func(i, j int) bool {
			return packingAsset[i].SellingStep > packingAsset[j].SellingStep
		})

		// Convert the timestamp to string in the desired format
		timestampStr := time.Unix(record.Timestamp.Seconds, int64(record.Timestamp.Nanos)).Format(utils.TIMEFORMAT)

		historyRecord := &models.PackingTransactionHistory{
			TxId:      record.TxId,
			Value:     packingAsset,
			Timestamp: timestampStr,
			IsDelete:  record.IsDelete,
		}

		latestHistory = historyRecord
	}

	if latestHistory == nil {
		return nil, fmt.Errorf("no history found for key %s", key)
	}

	return latestHistory, nil
}

func (s *SmartContract) GetHistoryForKey(ctx contractapi.TransactionContextInterface, key string) ([]*models.PackingTransactionHistory, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for key %s: %v", key, err)
	}
	defer resultsIterator.Close()

	var history []*models.PackingTransactionHistory
	var assetsValue []*models.PackingTransactionResponse

	for resultsIterator.HasNext() {
		// Get the next history record
		record, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next history record for key %s: %v", key, err)
		}

		var asset models.PackingTransactionResponse
		if !record.IsDelete {
			err = json.Unmarshal(record.Value, &asset)
			if err != nil {
				return nil, err
			}
			assetsValue = append(assetsValue, &asset)

		} else {
			assetsValue = []*models.PackingTransactionResponse{}
		}
		// Convert the timestamp to string in the desired format
		timestampStr := time.Unix(record.Timestamp.Seconds, int64(record.Timestamp.Nanos)).Format(utils.TIMEFORMAT)

		historyRecord := &models.PackingTransactionHistory{
			TxId:      record.TxId,
			Value:     assetsValue,
			Timestamp: timestampStr,
			IsDelete:  record.IsDelete,
		}

		history = append(history, historyRecord)
	}

	return history, nil
}