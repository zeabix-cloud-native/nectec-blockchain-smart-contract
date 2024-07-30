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
		CancelReason:   input.CancelReason,
		PackerId:       input.PackerId,
		Province:       input.Province,
		District:       input.District,
		Gmp:            input.Gmp,
		Gap:            input.Gap,
		ProcessStatus:  input.ProcessStatus,
		SellingStep:  	input.SellingStep,
		Owner:          clientIDPacking,
		OrgName:        orgName,
		UpdatedAt:      input.UpdatedAt,
		CreatedAt:      input.CreatedAt,
		DocType: 		models.Packing,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdatePacking(ctx contractapi.TransactionContextInterface, args string) error {
	fmt.Println("Start update packing")

	// Initialize entityPacking and unmarshal the input args into it
	var entityPacking models.TransactionPacking
	err := json.Unmarshal([]byte(args), &entityPacking)
	if err != nil {
		return fmt.Errorf("failed to unmarshal input: %v", err)
	}

	fmt.Println("Start read packing")
	asset, err := s.ReadPacking(ctx, entityPacking.Id)
	if err != nil {
		return fmt.Errorf("failed to read asset: %v", err)
	}

	fmt.Println("Modify packing model")
	asset.ForecastWeight = entityPacking.ForecastWeight
	asset.ActualWeight = entityPacking.ActualWeight
	asset.SavedTime = entityPacking.SavedTime
	asset.ApprovedDate = entityPacking.ApprovedDate
	asset.ApprovedType = entityPacking.ApprovedType
	asset.FinalWeight = entityPacking.FinalWeight
	asset.Province = entityPacking.Province
	asset.District = entityPacking.District
	asset.Remark = entityPacking.Remark
	asset.CancelReason = entityPacking.CancelReason
	asset.Gmp = entityPacking.Gmp
	asset.Gap = entityPacking.Gap
	asset.ProcessStatus = entityPacking.ProcessStatus
	asset.SellingStep = entityPacking.SellingStep
	asset.UpdatedAt = entityPacking.UpdatedAt

	fmt.Println("Modify packing done")
	assetJSON, errPacking := json.Marshal(asset)
	if errPacking != nil {
		return fmt.Errorf("failed to marshal updated asset: %v", errPacking)
	}

	fmt.Println("Send UpdateAsset event")
	err = ctx.GetStub().SetEvent("UpdateAsset", assetJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	fmt.Printf("Packing %s updated successfully\n", entityPacking.Id)
	return ctx.GetStub().PutState(entityPacking.Id, assetJSON)
}


func (s *SmartContract) DeletePacking(ctx contractapi.TransactionContextInterface, id string) error {

	assetPacking, err := s.ReadPacking(ctx, id)
	if err != nil {
		return err
	}

	// clientIDPacking, err := utils.GetIdentity(ctx)
	// utils.HandleError(err)

	// if clientIDPacking != assetPacking.Owner {
	// 	return utils.ReturnError(utils.UNAUTHORIZE)
	// }

	return ctx.GetStub().DelState(assetPacking.Id)
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

	// Build query string
	selector := map[string]interface{}{
		"selector": filterPacking,
	}

	if inputPacking.Search != nil && *inputPacking.Search != "" {
		searchTerm := *inputPacking.Search

		_, searchExists := filterPacking["search"]
		if searchExists {
			delete(filterPacking, "search")
		}

		selector["selector"] = map[string]interface{}{
			"$and": []map[string]interface{}{
				filterPacking,
				{
					"$or": []map[string]interface{}{
						{"gmp": map[string]interface{}{"$regex": searchTerm}},
						{"packingHouseName": map[string]interface{}{"$regex": searchTerm}},
						{"gap": map[string]interface{}{"$regex": searchTerm}},
					},
				},
			},
		}
	}

	queryString, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	// Debugging: print the query string to ensure it matches the expected filter criteria
	fmt.Printf("Query String for Fetching and Counting: %s\n", queryString)

    total, err := utils.CountTotalResults(ctx, string(queryString))
    if err != nil {
        return nil, err
    }

    arrPacking, _, err := utils.PackingFetchResultsWithPagination(ctx, inputPacking, selector)
    if err != nil {
        return nil, err
    }

    if len(arrPacking) == 0 {
        arrPacking = []*models.PackingTransactionResponse{}
    }

    // CalculateTotalSold(arrPacking)

    return &models.PackingGetAllResponse{
        Data:  "All Packing",
        Obj:   arrPacking,
        Total: total,
    }, nil
}

// Utility function to calculate totalSold
// func CalculateTotalSold(documents []*models.PackingTransactionResponse) {
//     gapTotals := make(map[string]float32)
    
//     // Calculate total sold for each gap
//     for _, doc := range documents {
//         if doc.Gap != "" {
//             gapTotals[doc.Gap] += doc.FinalWeight
//         }
//     }
    
//     // Update each document with the total sold value
//     for _, doc := range documents {
//         if doc.Gap != "" {
//             doc.TotalSold = gapTotals[doc.Gap]
//         }
//     }
// }

func CalculateTotalPackingSold(documents []*models.TransactionPacking) {
	packingTotals := make(map[string]float32)
	for _, doc := range documents {
		if doc.Gap != "" {
			packingTotals[doc.Gap] += doc.FinalWeight
		}
	}
	for _, doc := range documents {
		if doc.Gap != "" {
			doc.TotalSold = packingTotals[doc.Gap]
		}
	}
}

func (s *SmartContract) CalculateTotalSold(ctx contractapi.TransactionContextInterface, gapId string) (int, error) {
	packings, err := FetchPackingDocsByGap(ctx, gapId)
	if err != nil {
		return 0, err
	}

	totalSold := 0

	for _, packing := range packings {
		if (packing.ProcessStatus == 2 || packing.ProcessStatus == 3) {
			totalSold += int(packing.ActualWeight)
		}
	}

	return totalSold, nil
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

		// Check if docType is "packing"
		if docType, ok := m["docType"]; ok && docType == "packing" {
			// Apply the key-value filter
			if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
				assetPacking = append(assetPacking, &asset)
			}
		}
	}

	CalculateTotalPackingSold(assetPacking)

	sort.Slice(assetPacking, func(i, j int) bool {
        t1, err1 := time.Parse(time.RFC3339, assetPacking[i].CreatedAt)
        t2, err2 := time.Parse(time.RFC3339, assetPacking[j].CreatedAt)
        if err1 != nil || err2 != nil {
            fmt.Println("Error parsing time:", err1, err2)
            return false
        }
        return t1.After(t2)
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