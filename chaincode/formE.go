package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreateFormE(ctx contractapi.TransactionContextInterface, id string, formEJSON string) error {
	var formE models.TransactionFormE

	err := json.Unmarshal([]byte(formEJSON), &formE)
	if err != nil {
		return err
	}

    orgName, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
    }

    existPacker, err := utils.AssetExists(ctx, id)
    if err != nil {
        return fmt.Errorf("error checking if asset exists: %v", err)
    }
    if existPacker {
        return fmt.Errorf("the asset %s already exists", id)
    }

    clientID, err := utils.GetIdentity(ctx)
    if err != nil {
        return fmt.Errorf("failed to get submitting client's identity: %v", err)
    }

	formE.Id = id
    formE.DocType = models.FormE
    formE.Owner =   clientID
    formE.OrgName =   orgName

	formEAsBytes, err := json.Marshal(formE)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, formEAsBytes)
}

func (s *SmartContract) QueryFormEWithPagination(ctx contractapi.TransactionContextInterface, filterParams string) (*models.TransactionFormEResponse, error) {
	var filters models.FormEFilterParams
	err := json.Unmarshal([]byte(filterParams), &filters)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter parameters: %v", err)
	}

	selector := map[string]interface{}{
		"docType": "formE",
	}

	if filters.CreatedById != "" {
		selector["createdById"] = filters.CreatedById
	}

	if filters.ReferenceNo != "" {
		selector["referenceNo"] = filters.ReferenceNo
	}

	if filters.ExportNumber != "" {
		selector["invoice.exportNumber"] = filters.ExportNumber
	}

	if filters.Status != "" {
		selector["status"] = filters.Status
	}

	if filters.RequestType != "" {
		selector["requestType"] = filters.RequestType
	}

	if filters.StartDate != "" && filters.EndDate != "" {
		selector["createdAt"] = map[string]interface{}{
			"$gte": filters.StartDate,
			"$lte": filters.EndDate,
		}
	}

	if filters.Search != "" {
		elemMatch := map[string]interface{}{
			"$or": []map[string]interface{}{
				{"containerNumber": filters.Search},
				{"palletNumber": filters.Search},
			},
		}
		selector["$or"] = []map[string]interface{}{
			{"invoice.productAndPackaging": map[string]interface{}{"$elemMatch": elemMatch}},
			{"invoice.invoiceNumber": filters.Search},
		}
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

	// Log the count query string
	fmt.Printf("Count query string: %s\n", string(countQueryString))

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
		return &models.TransactionFormEResponse{
			Data:  []*models.TransactionFormE{},
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

	if err != nil {
		return nil, fmt.Errorf("failed to marshal query string: %v", err)
	}

	// Log the paginated query string
	fmt.Printf("Paginated query string: %s\n", string(queryString))

	// Execute the paginated query
	resultsIterator, _, err := ctx.GetStub().GetQueryResultWithPagination(string(queryString), int32(filters.Limit), "")
	if err != nil {
		return nil, fmt.Errorf("failed to get query result with pagination: %v", err)
	}
	defer resultsIterator.Close()

	var assets []*models.TransactionFormE
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var asset models.TransactionFormE
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		assets = append(assets, &asset)
	}

	return &models.TransactionFormEResponse{
		Data:  assets,
		Total: totalCount,
	}, nil
}



func (s *SmartContract) ReadFormE(ctx contractapi.TransactionContextInterface, referenceNo string) (*models.TransactionFormE, error) {
	formEAsBytes, err := ctx.GetStub().GetState(referenceNo)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if formEAsBytes == nil {
		return nil, fmt.Errorf("the formE %s does not exist", referenceNo)
	}

	var formE models.TransactionFormE
	err = json.Unmarshal(formEAsBytes, &formE)
	if err != nil {
		return nil, err
	}

	return &formE, nil
}

func (s *SmartContract) GetFormEHistoryForKey(ctx contractapi.TransactionContextInterface, key string) ([]*models.FormETransactionHistory, error) {
	formEResultsIterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for key %s: %v", key, err)
	}
	defer formEResultsIterator.Close()

	var formEHistory []*models.FormETransactionHistory

	for formEResultsIterator.HasNext() {
		// Get the next history record
		formE, err := formEResultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next history formE for key %s: %v", key, err)
		}

		var asset models.TransactionFormE
		if !formE.IsDelete {
			err = json.Unmarshal(formE.Value, &asset)
			if err != nil {
				return nil, err
			}

			// Ensure productAndPackaging is always initialized as an array
			if asset.Invoice.ProductAndPackaging == nil {
				asset.Invoice.ProductAndPackaging = []models.ProductAndPackaging{}
			}
		} else {
			// If the record is deleted, initialize asset to an empty array
			asset.Invoice.ProductAndPackaging = []models.ProductAndPackaging{}
		}

		var formEAssets []*models.TransactionFormE
		formEAssets = append(formEAssets, &asset)

		// Convert the timestamp to string in the desired format
		timestampStr := time.Unix(formE.Timestamp.Seconds, int64(formE.Timestamp.Nanos)).Format(utils.TIMEFORMAT)

		historyFormE := &models.FormETransactionHistory{
			TxId:      formE.TxId,
			Value:     formEAssets,
			Timestamp: timestampStr,
			IsDelete:  formE.IsDelete,
		}

		formEHistory = append(formEHistory, historyFormE)
	}

	return formEHistory, nil
}

func (s *SmartContract) CancelFormE(ctx contractapi.TransactionContextInterface,
	args string) error {
	entityType := models.TransactionFormE{}
	inputInterface, err := utils.Unmarshal(args, entityType)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionFormE)

	asset, err := s.ReadFormE(ctx, input.Id)
	utils.HandleError(err)

	asset.Id = input.Id
	asset.CancelReason = input.CancelReason
	asset.UpdatedAt = input.UpdatedAt
	asset.Status = "2"

	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) GetFormEByReferenceId(ctx contractapi.TransactionContextInterface, referenceId string) (*models.TransactionFormE, error) {
	queryFormE := fmt.Sprintf(`{"selector":{"referenceNo":"%s", "docType": "formE"}}`, referenceId)

	resultsFormE, err := ctx.GetStub().GetQueryResult(queryFormE)
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsFormE.Close()

	if !resultsFormE.HasNext() {
		return &models.TransactionFormE{}, nil
	}

	queryResponse, err := resultsFormE.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	var asset models.TransactionFormE
	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &asset, nil
}

