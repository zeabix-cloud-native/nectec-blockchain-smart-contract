package chaincode

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateExporter(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityExporter := models.TransactionExporter{}
	inputInterface, err := utils.Unmarshal(args, entityExporter)
	if err != nil {
		return err
	}
	input := inputInterface.(*models.TransactionExporter)

	// err := ctx.GetClientIdentity().AssertAttributeValue("exporter.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have exporter.creator role")
	}

	existExporter, err := utils.AssetExists(ctx, input.Id)
	if err != nil {
		return err
	}
	if existExporter {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	CreatedTime := utils.GetTimeNow()

	asset := models.TransactionExporter{
		Id:        input.Id,
		CertId:    input.CertId,
		PlantType:    input.PlantType,
		Name:    input.Name,
		Address:    input.Address,
		District:    input.District,
		Province:    input.Province,
		PostCode:    input.PostCode,
		Email:    input.Email,
		IssueDate:    input.IssueDate,
		ExpiredDate:    input.ExpiredDate,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: CreatedTime,
		CreatedAt: CreatedTime,
		DocType: models.Exporter,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) GetLastIdExporter(ctx contractapi.TransactionContextInterface) string {
	// Query to get all records sorted by ID in descending order
	query := `{
		"selector": {
			"docType": "exporter"
		},
		"sort": [{"_id": "desc"}],
		"limit": 1
	}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return fmt.Sprintf("error querying CouchDB: %v", err)
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

func (s *SmartContract) CreateExporterCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.TransactionExporter 
	var eventPayloads []models.TransactionExporter

	// Unmarshal input arguments
	errInputExporter := json.Unmarshal([]byte(args), &inputs)
	if errInputExporter != nil {
		return fmt.Errorf("failed to unmarshal input arguments: %v", errInputExporter)
	}

	// Process each input
	for _, input := range inputs {
		// Get client's MSP ID
		orgNameG, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		// Check if the asset already exists
		existExporter, err := utils.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existExporter {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		// Get client's identity
		clientIDG, err := utils.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		// Create the asset
		assetG := models.TransactionExporter{
			Id:          input.Id,
			CertId:      input.CertId,
			PlantType:   input.PlantType,
			Name:        input.Name,
			Address:     input.Address,
			District:    input.District,
			Province:    input.Province,
			PostCode:    input.PostCode,
			Email:       input.Email,
			IssueDate:   input.IssueDate,
			ExpiredDate: input.ExpiredDate,
			Owner:       clientIDG,
			OrgName:     orgNameG,
			DocType:     models.Exporter,
			CreatedAt:   utils.GetTimeNow(),
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

	// Marshal event payloads to JSON
	eventPayloadJSON, err := json.Marshal(eventPayloads)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload JSON: %v", err)
	}

	// Set the event
	ctx.GetStub().SetEvent("batchCreatedExporterEvent", eventPayloadJSON)

	return nil
}

func (s *SmartContract) UpdateExporter(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityExporter := models.TransactionExporter{}
	inputInterface, err := utils.Unmarshal(args, entityExporter)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionExporter)

	asset, err := s.ReadExporter(ctx, input.Id)
	utils.HandleError(err)

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientID != asset.Owner {
		return fmt.Errorf(utils.UNAUTHORIZE)
	}

	UpdatedTime := utils.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedTime
	asset.PlantType = input.PlantType
	asset.Name = input.Name
	asset.Address = input.Address
	asset.PostCode = input.PostCode
	asset.District = input.District
	asset.Province = input.Province
	asset.Email = input.Email
	asset.IssueDate = input.IssueDate
	asset.ExpiredDate = input.ExpiredDate
	asset.UpdatedAt = utils.GetTimeNow()

	assetJSON, errE := json.Marshal(asset)
	utils.HandleError(errE)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteExporter(ctx contractapi.TransactionContextInterface, id string) error {

	assetE, err := s.ReadExporter(ctx, id)
	utils.HandleError(err)

	// clientIDExporter, err := utils.GetIdentity(ctx)
	// utils.HandleError(err)

	// if clientIDExporter != assetE.Owner {
	// 	return fmt.Errorf(utils.UNAUTHORIZE)
	// }

	return ctx.GetStub().DelState(assetE.Id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetE, err := s.ReadExporter(ctx, id)
	utils.HandleError(err)

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientID != assetE.Owner {
		return fmt.Errorf(utils.UNAUTHORIZE)
	}

	assetE.Owner = newOwner
	assetJSON, err := json.Marshal(assetE)
	utils.HandleError(err)
	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadExporter(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionExporter, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionExporter
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) GetAllExporter(ctx contractapi.TransactionContextInterface, args string) (*models.ExporterGetAllResponse, error) {
	entityGetAll := models.ExporterFilterGetAll{}
	interfaceE, err := utils.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	inputExporter := interfaceE.(*models.ExporterFilterGetAll)
	filterExporter := utils.ExporterSetFilter(inputExporter)

	arrExporter, total, err := utils.ExporterFetchResultsWithPagination(ctx, inputExporter, filterExporter)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrExporter, func(i, j int) bool {
		return arrExporter[i].CreatedAt.After(arrExporter[j].CreatedAt)
	})

	if len(arrExporter) == 0 {
		arrExporter = []*models.ExporterTransactionResponse{}
	}

	return &models.ExporterGetAllResponse{
		Data:  "All Exporter",
		Obj:   arrExporter,
		Total: total,
	}, nil
}


func (s *SmartContract) FilterExporter(ctx contractapi.TransactionContextInterface, key, value string) ([]*models.TransactionExporter, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetExporter []*models.TransactionExporter
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.TransactionExporter
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetExporter = append(assetExporter, &asset)
		}
	}

	sort.Slice(assetExporter, func(i, j int) bool {
		return assetExporter[i].UpdatedAt.After(assetExporter[j].UpdatedAt)
	})

	return assetExporter, nil
}