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


func (s *SmartContract) CreatePacker(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityPacker := models.TransactionPacker{}
	inputInterface, err := utils.Unmarshal(args, entityPacker)
	if err != nil {
		return err
	}
	input := inputInterface.(*models.TransactionPacker)

	// err := ctx.GetClientIdentity().AssertAttributeValue("packer.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have packer.creator role")
	}

	existPacker, err := utils.AssetExists(ctx, input.Id)
	utils.HandleError(err)
	if existPacker {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	// timestamp := utils.GenerateTimestamp()

	asset := models.TransactionPacker{
		Id:        input.Id,
		CertId:    input.CertId,
		UserId:    input.UserId,
		IsCanExport: input.IsCanExport,
		PackingHouseName: input.PackingHouseName,
		PackingHouseRegisterNumber: input.PackingHouseRegisterNumber,
		Owner:     clientID,
		OrgName:   orgName,
		DocType:   models.Packer,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdatePacker(ctx contractapi.TransactionContextInterface, args string) error {
    entityPacker := models.TransactionPacker{}
    inputInterface, err := utils.Unmarshal(args, entityPacker)
    if err != nil {
        return err
    }
    input := inputInterface.(*models.TransactionPacker)

    asset, err := s.ReadPacker(ctx, input.Id)
    if err != nil {
        return err
    }
	
    // clientID, err := utils.GetIdentity(ctx)
    // if err != nil {
    //     return err
    // }

    // if clientID != asset.Owner {
    //     return fmt.Errorf(utils.UNAUTHORIZE)
    // }

	asset.CertId = input.CertId
	asset.UserId = input.UserId
	asset.PackerGmp = nil
	asset.PackingHouseName = input.PackingHouseName
	asset.PackingHouseRegisterNumber = input.PackingHouseRegisterNumber
	asset.IsCanExport = input.IsCanExport
	

	timestamp := utils.GenerateTimestamp()
	asset.UpdatedAt = timestamp

    assetJSON, errP := json.Marshal(asset)
    if errP != nil {
        return errP
    }

    return ctx.GetStub().PutState(input.Id, assetJSON)
}


func (s *SmartContract) DeletePacker(ctx contractapi.TransactionContextInterface, id string) error {

	assetPacker, err := s.ReadPacker(ctx, id)
	utils.HandleError(err)

	// clientIDPacker, err := utils.GetIdentity(ctx)
	// utils.HandleError(err)

	// if clientIDPacker != assetPacker.Owner {
	// 	return fmt.Errorf(utils.UNAUTHORIZE)
	// }

	return ctx.GetStub().DelState(assetPacker.Id)
}

func (s *SmartContract) GetPackerByPackerId(ctx contractapi.TransactionContextInterface, packerId string) (*models.PackerByIdResponse, error) {
	queryKeyFarmer := fmt.Sprintf(`{"selector":{"userId":"%s", "docType": "packer"}}`, packerId)

	resultsIteratorFarmer, err := ctx.GetStub().GetQueryResult(queryKeyFarmer)
	var asset *models.PackerTransactionResponse
	resData := "Get packer by packerId"
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIteratorFarmer.Close()

	if !resultsIteratorFarmer.HasNext() {
		resData = "Not found packer by packerId"

		return &models.PackerByIdResponse{
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

	// Attach related GMP documents
	queryString := fmt.Sprintf(`{
		"selector": {
			"docType": "gmp",
			"packerId": "%s"
		}
	}`, asset.Id)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query related gmp documents: %v", err)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var gmpDoc models.PackerGmp
		err = json.Unmarshal(queryResponse.Value, &gmpDoc)
		if err != nil {
			return nil, err
		}

		asset.PackerGmp = gmpDoc
	}

	return &models.PackerByIdResponse{
		Data: resData,
		Obj:  asset,
	}, nil
}

func (s *SmartContract) ReadPacker(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionPacker, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionPacker
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	// Attach related GMP documents
	queryString := fmt.Sprintf(`{
		"selector": {
			"docType": "gmp",
			"packerId": "%s"
		}
	}`, asset.Id)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query related gmp documents: %v", err)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var gmpDoc *models.PackerGmp
		err = json.Unmarshal(queryResponse.Value, &gmpDoc)
		if err != nil {
			return nil, err
		}

		asset.PackerGmp = gmpDoc
	}

	return &asset, nil
}

func (s *SmartContract) GetPackerById(ctx contractapi.TransactionContextInterface, id string) (*models.PackerTransactionResponse, error) {
	queryPacker := fmt.Sprintf(`{"selector":{"id":"%s"}}`, id)

	resultsPacker, err := ctx.GetStub().GetQueryResult(queryPacker)
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsPacker.Close()

	if !resultsPacker.HasNext() {
		return nil, fmt.Errorf("the asset with id %s does not exist", id)
	}

	queryResponse, err := resultsPacker.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	var asset models.PackerTransactionResponse
	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &asset, nil
}

func (s *SmartContract) GetAllPacker(ctx contractapi.TransactionContextInterface, args string) (*models.PackerGetAllResponse, error) {
    var filterPacker = map[string]interface{}{}
    filterPacker["docType"] = "packer"

    entityGetAll := models.FilterGetAllPacker{}
    interfacePacker, err := utils.Unmarshal(args, entityGetAll)
    if err != nil {
        return nil, err
    }
    input := interfacePacker.(*models.FilterGetAllPacker)

    arrPacker, total, err := utils.PackerFetchResultsWithPagination(ctx, input, filterPacker)
    if err != nil {
        return nil, err
    }

	sort.Slice(arrPacker, func(i, j int) bool {
        t1, err1 := time.Parse(time.RFC3339, arrPacker[i].CreatedAt)
        t2, err2 := time.Parse(time.RFC3339, arrPacker[j].CreatedAt)
        if err1 != nil || err2 != nil {
            fmt.Println("Error parsing time:", err1, err2)
            return false
        }
        return t1.After(t2)
    })

    // Attach related GMP documents
    for _, packer := range arrPacker {
        queryString := fmt.Sprintf(`{
            "selector": {
                "docType": "gmp",
                "packerId": "%s"
            }
        }`, packer.Id)

        resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
        if err != nil {
            return nil, fmt.Errorf("failed to query related gmp documents: %v", err)
        }
        defer resultsIterator.Close()

        for resultsIterator.HasNext() {
            queryResponse, err := resultsIterator.Next()
            if err != nil {
                return nil, err
            }

            var gmpDoc models.PackerGmp
            err = json.Unmarshal(queryResponse.Value, &gmpDoc)
            if err != nil {
                return nil, err
            }

            packer.PackerGmp = gmpDoc
        }
    }

    if len(arrPacker) == 0 {
        arrPacker = []*models.PackerTransactionResponse{}
    }

    return &models.PackerGetAllResponse{
        Data:  "All Packer",
        Obj:   arrPacker,
        Total: total,
    }, nil
}

func (s *SmartContract) GetLastIdPacker(ctx contractapi.TransactionContextInterface) string {
	// Query to get all records sorted by ID in descending order
	query := `{
		"selector": {
			"docType": "packer"
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

func (s *SmartContract) CreatePackerCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []models.TransactionPacker
	var eventPayloads []models.TransactionPacker

	errPackerInput := json.Unmarshal([]byte(args), &inputs)
	if errPackerInput != nil {
		return fmt.Errorf("failed to unmarshal JSON array: %v", errPackerInput)
	}

	for _, input := range inputs {
		fmt.Printf("create packer csv input %v", input)
		orgName, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		existPacker, err := utils.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existPacker {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		clientID, err := utils.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		// timestamp := utils.GenerateTimestamp()

		asset := models.TransactionPacker{
			Id:        input.Id,
			CertId:    input.CertId,
			UserId:    input.UserId,
			IsCanExport: input.IsCanExport,
			PackingHouseName: input.PackingHouseName,
			PackingHouseRegisterNumber: input.PackingHouseRegisterNumber,
			Owner:     clientID,
			OrgName:   orgName,
			UpdatedAt: input.UpdatedAt,
			CreatedAt: input.CreatedAt,
			DocType:   models.Packer,
		}

		packerAssetJSON, packerErr := json.Marshal(asset)
		eventPayloads = append(eventPayloads, asset)
		if packerErr != nil {
			return fmt.Errorf("failed to marshal asset JSON: %v", err)
		}

		err = ctx.GetStub().PutState(input.Id, packerAssetJSON)
		if err != nil {
			return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
		}

		fmt.Printf("Asset %s created successfully\n", input.Id)
	}

	eventPayloadJSON, err := json.Marshal(eventPayloads)
	if err != nil {
		return fmt.Errorf("failed to marshal asset JSON: %v", err)
	}

	ctx.GetStub().SetEvent("batchCreatedPackerEvent", eventPayloadJSON)

	return nil
}