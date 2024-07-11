package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreatePackagingCsv(
    ctx contractapi.TransactionContextInterface,
    args string,
) error {
    var inputs []models.TransactionPackaging

    errInputPackaging := json.Unmarshal([]byte(args), &inputs)
    if errInputPackaging != nil {
        return fmt.Errorf("error unmarshaling input: %v", errInputPackaging)
    }

    var batchErrors []error

    for _, input := range inputs {
        if err := s.processSinglePackaging(ctx, input); err != nil {
            batchErrors = append(batchErrors, err)
            fmt.Printf("Error processing asset %s: %v\n", input.Id, err)
        } else {
            fmt.Printf("Packaging Asset %s created successfully\n", input.Id)
        }
    }

    if len(batchErrors) > 0 {
        return fmt.Errorf("encountered errors during processing: %v", batchErrors)
    }

    return nil
}

func (s *SmartContract) processSinglePackaging(ctx contractapi.TransactionContextInterface, input models.TransactionPackaging) error {
    orgNamePackaging, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
    }

    existPackaging, err := utils.AssetExists(ctx, input.Id)
    if err != nil {
        return fmt.Errorf("error checking if asset exists: %v", err)
    }
    if existPackaging {
        return fmt.Errorf("the asset %s already exists", input.Id)
    }

    clientIDPackaging, err := utils.GetIdentity(ctx)
    if err != nil {
        return fmt.Errorf("failed to get submitting client's identity: %v", err)
    }

    // timestamp := utils.GenerateTimestamp()

    assetPackaging := models.TransactionPackaging{
        Id:          input.Id,
        ContainerId: input.ContainerId,
        ExportId:    input.ExportId,
        LotNumber:   input.LotNumber,
        BoxId:       input.BoxId,
        Gap:         input.Gap,
        Gmp:         input.Gmp,
        GradeName:   input.GradeName,
        Gtin14:      input.Gtin14,
        Gtin13:      input.Gtin13,
        VarietyName: input.VarietyName,
        DocType:     models.Packaging,
        Owner:       clientIDPackaging,
        OrgName:     orgNamePackaging,
        CreatedAt:   input.CreatedAt,
        UpdatedAt:   input.UpdatedAt,
    }

    assetJSON, err := json.Marshal(assetPackaging)
    if err != nil {
        return fmt.Errorf("failed to marshal asset JSON for asset %s: %v", input.Id, err)
    }

    err = ctx.GetStub().PutState(input.Id, assetJSON)
    if err != nil {
        return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
    }

    return nil
}
