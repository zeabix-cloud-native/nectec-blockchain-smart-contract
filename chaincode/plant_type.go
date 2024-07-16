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