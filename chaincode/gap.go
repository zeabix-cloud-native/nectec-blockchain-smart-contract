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

func (s *SmartContract) CreateGAP(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityGap := models.TransactionGap{}
	inputInterface, err := utils.Unmarshal(args, entityGap)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionGap)

	// err := ctx.GetClientIdentity().AssertAttributeValue("gap.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have gap.creator role1")
	}

	existsGap, err := utils.AssetExists(ctx, input.Id)
	utils.HandleError(err)
	if existsGap {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientIDGap, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	asset := models.TransactionGap{
		Id:          				input.Id,
		CertID:      				input.CertID,
		DisplayCertID:      input.DisplayCertID,
		AreaCode:    input.AreaCode,
		AreaRai:     input.AreaRai,
		AreaStatus:  input.AreaStatus,
		OldAreaCode: input.OldAreaCode,
		IssueDate:   input.IssueDate,
		ExpireDate:  input.ExpireDate,
		District:    input.District,
		Province:    input.Province,
		UpdatedAt:   input.UpdatedAt,
		Source:      input.Source,
		FarmerID:    input.FarmerID,
		Owner:       clientIDGap,
		OrgName:     orgName,
		CreatedAt:   input.CreatedAt,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateGap(ctx contractapi.TransactionContextInterface, args string) error {

	entityGap := models.TransactionGap{}
	inputInterface, err := utils.Unmarshal(args, entityGap)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionGap)

	asset, err := s.ReadGap(ctx, input.Id)
	utils.HandleError(err)

	asset.Id = input.Id
	asset.DisplayCertID = input.DisplayCertID
	asset.CertID = input.CertID
	asset.AreaCode = input.AreaCode
	asset.AreaRai = input.AreaRai
	asset.AreaStatus = input.AreaStatus
	asset.OldAreaCode = input.OldAreaCode
	asset.IssueDate = input.IssueDate
	asset.ExpireDate = input.ExpireDate
	asset.District = input.District
	asset.Province = input.Province
	asset.Source = input.Source
	asset.FarmerID = input.FarmerID
	asset.UpdatedAt = input.UpdatedAt

	assetJSON, errGap := json.Marshal(asset)
	utils.HandleError(errGap)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) ReadGap(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionGap, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionGap
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	asset.IsCanDelete = true

	salesQueryString := fmt.Sprintf(`{
		"selector": {
			"docType": "packing",
			"gap": "%s"
		},
		"use_index": ["_design/index-DocTypeGap", "index-DocTypeGap"]
	}`, asset.CertID)

	fmt.Printf("salesQueryString %v", salesQueryString)

	salesResultsIterator, err := ctx.GetStub().GetQueryResult(salesQueryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query related sales: %v", err)
	}
	defer salesResultsIterator.Close()

	if salesResultsIterator.HasNext() {
		asset.IsCanDelete = false
	}

	return &asset, nil
}


func (s *SmartContract) DeleteGap(ctx contractapi.TransactionContextInterface, id string) error {
	assetGap, err := s.ReadGap(ctx, id)
	utils.HandleError(err)


	if (assetGap.FarmerID != "") {
		queryFarmer := fmt.Sprintf(`{"selector":{"docType":"farmer","id":"%s"}}`, assetGap.FarmerID)
		fmt.Printf("[DeleteGap] queryFarmer %v \n", queryFarmer);
		resultsIteratorFarmer, err := ctx.GetStub().GetQueryResult(queryFarmer)
		if err != nil {
			return fmt.Errorf("error querying chaincode: %v", err)
		}
		defer resultsIteratorFarmer.Close()

		if resultsIteratorFarmer.HasNext() {
			queryResponse, err := resultsIteratorFarmer.Next()
			if err != nil {
				return fmt.Errorf("error iterating query results: %v", err)
			}

			var farmer *models.TransactionFarmer
			err = json.Unmarshal(queryResponse.Value, &farmer)
			if err != nil {
				return fmt.Errorf("error unmarshalling JSON: %v", err)
			}
			fmt.Printf("[DeleteGap] farmer %v \n", farmer);

			var updatedFarmerGaps []models.FarmerGap

			if farmer.FarmerGaps != nil && len(farmer.FarmerGaps) > 0 {
				for _, item := range farmer.FarmerGaps {
					if item.CertID != assetGap.CertID {
						updatedFarmerGaps = append(updatedFarmerGaps, item)
					}
				}
				farmer.FarmerGaps = updatedFarmerGaps
			} else {
				farmer.FarmerGaps = []models.FarmerGap{}
			}
			
			if (farmer.FarmerGaps != nil && len(farmer.FarmerGaps) > 0) {
				for _, item := range farmer.FarmerGaps {
					if (item.CertID != assetGap.CertID) {
						farmer.FarmerGaps = append(farmer.FarmerGaps, item);
					}
				}
			} 

			if len(updatedFarmerGaps) == 0 {
				fmt.Println("[DeleteGap] No gaps remaining after deletion")
			}

			farmer.FarmerGaps = updatedFarmerGaps

			assetJSON, err := json.Marshal(farmer);
			fmt.Printf("[DeleteGap] assetJSON %v \n", farmer);
			if err != nil {
				return fmt.Errorf("error marshalling JSON: %v", err)
			}
			ctx.GetStub().PutState(farmer.Id, assetJSON);
			fmt.Printf("Update farmer from delete gap successfully \n");
		}
	}

	return ctx.GetStub().DelState(assetGap.Id)
}

func (s *SmartContract) GetGapByFarmerID(ctx contractapi.TransactionContextInterface, farmerId string) (*models.GetGapByCertIdResponse, error) {
	// Get the asset using farmerId 
	queryKeyFarmer := fmt.Sprintf(`{"selector":{"farmerId":"%s"}}`, farmerId)

	resultsIteratorFarmer, err := ctx.GetStub().GetQueryResult(queryKeyFarmer)
	var asset *models.GapTransactionResponse
	resData := "Get gap by farmerId"
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIteratorFarmer.Close()

	if !resultsIteratorFarmer.HasNext() {
		resData = "Not found gap by farmerId"

		return &models.GetGapByCertIdResponse{
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

	return &models.GetGapByCertIdResponse{
		Data: resData,
		Obj:  asset,
	}, nil
}


func (s *SmartContract) GetGapByCertID(ctx contractapi.TransactionContextInterface, certID string) (*models.GetGapByCertIdResponse, error) {
	// Get the asset using CertID
	queryKeyGap := fmt.Sprintf(`{"selector":{"certId":"%s", "docType": "gap"}}`, certID)

	resultsIteratorGap, err := ctx.GetStub().GetQueryResult(queryKeyGap)
	var asset *models.GapTransactionResponse
	resData := "Get gap by certID"
	if err != nil {
		resData = fmt.Sprintf("error querying chaincode: %v", err)
        return &models.GetGapByCertIdResponse {
            Data: resData,
            Obj:  asset,
        }, nil
	}
	defer resultsIteratorGap.Close()

	if !resultsIteratorGap.HasNext() {
		resData = "Not found gap by certID"
        return &models.GetGapByCertIdResponse {
            Data: resData,
            Obj:  asset,
        }, nil
	}

	queryResponse, err := resultsIteratorGap.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &models.GetGapByCertIdResponse{
		Data: resData,
		Obj:  asset,
	}, nil
}

func (s *SmartContract) GetAllGAP(ctx contractapi.TransactionContextInterface, args string) (*models.GetAllGapResponse, error) {
	entityGetAllGap := models.FilterGetAllGap{}
	interfaceGap, err := utils.Unmarshal(args, entityGetAllGap)
	if err != nil {
		return nil, err
	}
	inputGap := interfaceGap.(*models.FilterGetAllGap)
	filterGap := utils.GapSetFilter(inputGap)

	// Build query string
	selector := map[string]interface{}{
		"selector": filterGap,
	}

	if inputGap.CertID != nil && *inputGap.CertID != "" {
		searchTerm := *inputGap.CertID
		selector["selector"] = map[string]interface{}{
			"$and": []map[string]interface{}{
				filterGap,
				{
					"$or": []map[string]interface{}{
						{"certId": map[string]interface{}{"$regex": searchTerm}},
						{"displayCertId": map[string]interface{}{"$regex": searchTerm}},
					},
				},
			},
		}
	}

	queryStringGap, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	// Debugging: print the query string to ensure it matches the expected filter criteria
	fmt.Printf("Query String for Fetching and Counting: %s\n", queryStringGap)

	// Count total results matching the query criteria
	total, err := utils.CountTotalResults(ctx, string(queryStringGap))
	if err != nil {
		return nil, err
	}

	// Fetch paginated results
	assets, err := utils.GapFetchResultsWithPagination(ctx, inputGap, filterGap)
	if err != nil {
		return nil, err
	}

	gapTotals := make(map[string]float32)
	for _, asset := range assets {
		packingDocs, err := FetchPackingDocsByGap(ctx, asset.CertID)
		if err != nil {
			return nil, err
		}

		for _, doc := range packingDocs {
			gapTotals[asset.CertID] += doc.FinalWeight
		}

		// Initialize isCanDelete to true
		asset.IsCanDelete = true

		// Query CouchDB for related 'packing' documents to determine if gap can be deleted
		salesQueryString := fmt.Sprintf(`{
			"selector": {
				"docType": "packing",
				"gap": "%s"
			}
		}`, asset.CertID)

		salesResultsIterator, err := ctx.GetStub().GetQueryResult(salesQueryString)
		if err != nil {
			return nil, fmt.Errorf("failed to query related sales: %v", err)
		}
		defer salesResultsIterator.Close()

		// If there are any related sales, set isCanDelete to false
		if salesResultsIterator.HasNext() {
			asset.IsCanDelete = false
		}
	}

	if len(assets) == 0 {
		assets = []*models.GapTransactionResponse{}
	}

	for _, asset := range assets {
		if total, ok := gapTotals[asset.CertID]; ok {
			asset.TotalSold = total
		}
	}

	return &models.GetAllGapResponse{
		Data:  "All Gap",
		Obj:   assets,
		Total: total,
	}, nil
}



func FetchPackingDocsByGap(ctx contractapi.TransactionContextInterface, gap string) ([]*models.PackingTransactionResponse, error) {
    filter := map[string]interface{}{
        "selector": map[string]interface{}{
            "docType": "packing",
            "gap":     gap,
        },
    }

    getStringPacking, err := json.Marshal(filter)
    if err != nil {
        return nil, err
    }

    queryPacking, err := ctx.GetStub().GetQueryResult(string(getStringPacking))
    if err != nil {
        return nil, err
    }
    defer queryPacking.Close()

    var dataPacking []*models.PackingTransactionResponse
    for queryPacking.HasNext() {
        queryResponse, err := queryPacking.Next()
        if err != nil {
            return nil, err
        }

        var asset models.PackingTransactionResponse
        err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
            return nil, err
        }

        dataPacking = append(dataPacking, &asset)
    }

    return dataPacking, nil
}

func (s *SmartContract) FilterGap(ctx contractapi.TransactionContextInterface, key, value string) ([]*models.TransactionGap, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetGap []*models.TransactionGap
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.TransactionGap
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetGap = append(assetGap, &asset)
		}
	}

	sort.Slice(assetGap, func(i, j int) bool {
        t1, err1 := time.Parse(time.RFC3339, assetGap[i].UpdatedAt)
        t2, err2 := time.Parse(time.RFC3339, assetGap[j].UpdatedAt)
        if err1 != nil || err2 != nil {
            fmt.Println("Error parsing time:", err1, err2)
            return false
        }
        return t1.After(t2)
    })

	return assetGap, nil
}

func (s *SmartContract) UpdateMultipleGap(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.TransactionGap

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
		
		var existingAsset models.TransactionGap
		err = json.Unmarshal(assetJSON, &existingAsset)
		if err != nil {
			return fmt.Errorf("failed to unmarshal existing asset: %v", err)
		}

		existingAsset.Id =          				 input.Id
		existingAsset.DisplayCertID =       input.DisplayCertID
		existingAsset.CertID =      input.CertID
		existingAsset.AreaCode =    input.AreaCode
		existingAsset.AreaRai =     input.AreaRai
		existingAsset.AreaStatus =  input.AreaStatus
		existingAsset.OldAreaCode = input.OldAreaCode
		existingAsset.IssueDate =   input.IssueDate
		existingAsset.ExpireDate =  input.ExpireDate
		existingAsset.District =    input.District
		existingAsset.Province =    input.Province
		existingAsset.Source =      input.Source
		existingAsset.FarmerID =    input.FarmerID
		existingAsset.UpdatedAt = 	input.UpdatedAt
		
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

func (s *SmartContract) CreateGapCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.TransactionGap

	errInputGap := json.Unmarshal([]byte(args), &inputs)
	utils.HandleError(errInputGap)

	for _, input := range inputs {
		// err := ctx.GetClientIdentity().AssertAttributeValue("gap.creator", "true")

		orgNameGap, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		existGap, err := utils.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existGap {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		clientIDGap, err := utils.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		assetGap := models.TransactionGap {
			Id:          				 input.Id,
			DisplayCertID:       input.DisplayCertID,
			CertID:      input.CertID,
			AreaCode:    input.AreaCode,
			AreaRai:     input.AreaRai,
			AreaStatus:  input.AreaStatus,
			OldAreaCode: input.OldAreaCode,
			IssueDate:   input.IssueDate,
			ExpireDate:  input.ExpireDate,
			District:    input.District,
			Province:    input.Province,
			UpdatedAt:   input.UpdatedAt,
			Source:      input.Source,
			FarmerID:    input.FarmerID,
			Owner:       clientIDGap,
			OrgName:     orgNameGap,
			DocType:     models.Gap,
			IsCanDelete: true,
			CreatedAt:   input.CreatedAt,
		}
		
		assetJSON, err := json.Marshal(assetGap)
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

