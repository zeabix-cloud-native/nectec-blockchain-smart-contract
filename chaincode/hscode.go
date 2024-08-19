package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	models "github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreateTransactionHscodes(ctx contractapi.TransactionContextInterface, hscodesJSON string) error {
	var inputs []models.TransactionHscode

    errInputHscode := json.Unmarshal([]byte(hscodesJSON), &inputs)
	utils.HandleError(errInputHscode)

    for _, input := range inputs {
        orgNameHscode, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		existHscode, err := utils.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existHscode {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		clientHscode, err := utils.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

        hscodeAsset := models.TransactionHscode{
			Id:                         input.Id,
            Hscode:                     input.Hscode,
            Description:                input.Description,
            OrgName:                    orgNameHscode,
			Order:                      input.Order,
			Owner:                      clientHscode,
			DocType:                    models.Hscode,
			CreatedAt:  				input.CreatedAt,		
			UpdatedAt:  				input.UpdatedAt,		
		}
        assetJSON, err := json.Marshal(hscodeAsset)
		if err != nil {
			return fmt.Errorf("failed to marshal asset JSON: %v", err)
		}

        err = ctx.GetStub().PutState(input.Id, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
		}

		fmt.Printf("Hscode Asset %s created successfully\n", input.Id)
    }

	return nil
}

func (s *SmartContract) DeleteAllHscodes(ctx contractapi.TransactionContextInterface) error {
	queryString := `{"selector":{"docType":"hscode"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return fmt.Errorf("failed to iterate query results: %v", err)
		}

		// Delete the asset
		err = ctx.GetStub().DelState(queryResponse.Key)
		if err != nil {
			return fmt.Errorf("failed to delete state for asset %s: %v", queryResponse.Key, err)
		}

		fmt.Printf("HS code asset %s deleted successfully\n", queryResponse.Key)
	}

	return nil
}

func (s *SmartContract) QueryHscodeWithPagination(ctx contractapi.TransactionContextInterface, filterParams string) (*models.TransactionHscodeResponse, error) {
	var filters models.PackagingFilterParams
	err := json.Unmarshal([]byte(filterParams), &filters)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal filter parameters: %v", err)
	}

	selector := map[string]interface{}{
		"docType": "hscode",
	}

	// Create query string for counting total records
	countQueryString, err := json.Marshal(map[string]interface{}{
		"selector": selector,
		"sort": []map[string]string{
			{"order": "asc"},
		},
		"use_index": []string{
            "_design/index-Order",
            "index-Order",
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
		return &models.TransactionHscodeResponse{
			Data:  []*models.TransactionHscode{},
			Total: totalCount,
		}, nil
	}

	// Create query string for paginated results
	queryString, err := json.Marshal(map[string]interface{}{
		"selector": selector,
		"sort": []map[string]string{
			{"order": "asc"},
		},
		"use_index": []string{
            "_design/index-Order",
            "index-Order",
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

	var assets []*models.TransactionHscode
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var asset models.TransactionHscode
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
		}
		assets = append(assets, &asset)
	}

	return &models.TransactionHscodeResponse{
		Data:  assets,
		Total: totalCount,
	}, nil
}