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

func (s *SmartContract) CreateNectecStaff(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityNstda := models.TransactionNectecStaff{}
	inputInterface, err := utils.Unmarshal(args, entityNstda)
	if err != nil {
		return err
	}
	input := inputInterface.(*models.TransactionNectecStaff)

	// err := ctx.GetClientIdentity().AssertAttributeValue("nstdaStaff.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have nstdaStaff.creator role")
	}

	existNstda, err := utils.AssetExists(ctx, input.Id)
	utils.HandleError(err)
	if existNstda {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	timestamp := utils.GenerateTimestamp()

	asset := models.TransactionNectecStaff{
		Id:        input.Id,
		CertId:    input.CertId,
		ProfileImg:    input.ProfileImg,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: timestamp,
		CreatedAt: timestamp,
		DocType: models.Nectec,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateNectecStaff(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityNstda := models.TransactionNectecStaff{}
	inputInterface, err := utils.Unmarshal(args, entityNstda)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionNectecStaff)

	asset, err := s.ReadNectecStaff(ctx, input.Id)
	utils.HandleError(err)

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientID != asset.Owner {
		return fmt.Errorf(utils.UNAUTHORIZE)
	}

	timestamp := utils.GenerateTimestamp()
	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = timestamp

	assetJSON, errN := json.Marshal(asset)
	utils.HandleError(errN)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteNectecStaff(ctx contractapi.TransactionContextInterface, id string) error {

	assetNstda, err := s.ReadNectecStaff(ctx, id)
	utils.HandleError(err)

	// clientIDNstda, err := utils.GetIdentity(ctx)
	// utils.HandleError(err)

	// if clientIDNstda != assetNstda.Owner {
	// 	return fmt.Errorf(utils.UNAUTHORIZE)
	// }

	return ctx.GetStub().DelState(assetNstda.Id)
}

func (s *SmartContract) DeleteNectecStaffFromCertId(ctx contractapi.TransactionContextInterface, certId string) error {
	queryString := fmt.Sprintf(`{"selector":{"certId":"%s"}}`, certId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return fmt.Errorf("failed to query certId: %v", err)
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext() {
		return fmt.Errorf("no asset found with certId: %s", certId)
	}

	// Assuming there is only one asset per certId. If there could be multiple, you'd need to handle that case.
	queryResponse, err := resultsIterator.Next()
	if err != nil {
		return fmt.Errorf("failed to get query response: %v", err)
	}

	// Convert queryResponse to asset structure if needed, e.g., NectecStaff
	var assetNstda models.TransactionNectecStaff
	err = json.Unmarshal(queryResponse.Value, &assetNstda)
	if err != nil {
		return fmt.Errorf("failed to unmarshal query response: %v", err)
	}

	return ctx.GetStub().DelState(assetNstda.Id)
}

func (s *SmartContract) QueryNectecStaffByCertId(ctx contractapi.TransactionContextInterface, certId string) ([]*models.TransactionNectecStaff, error) {
	queryString := fmt.Sprintf(`{"selector":{"certId":"%s","docType": "nectec}}`, certId)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query by certId: %v", err)
	}
	defer resultsIterator.Close()

	var assets []*models.TransactionNectecStaff
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate results: %v", err)
		}

		var asset models.TransactionNectecStaff
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *SmartContract) ReadNectecStaff(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionNectecStaff, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionNectecStaff
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}


func (s *SmartContract) GetAllNectecStaff(ctx contractapi.TransactionContextInterface, args string) (*models.GetAllNectecStaffResponse, error) {

	var filterNstda = map[string]interface{}{}

	entityGetAll := models.FilterGetAllNectecStaff{}
	interfaceNstda, err := utils.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := interfaceNstda.(*models.FilterGetAllNectecStaff)

	arrNstda, total, err := utils.NectecStaffFetchResultsWithPagination(ctx, input, filterNstda)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrNstda, func(i, j int) bool {
        t1, err1 := time.Parse(time.RFC3339, arrNstda[i].CreatedAt)
        t2, err2 := time.Parse(time.RFC3339, arrNstda[j].CreatedAt)
        if err1 != nil || err2 != nil {
            fmt.Println("Error parsing time:", err1, err2)
            return false
        }
        return t1.After(t2)
    })

	if len(arrNstda) == 0 {
		arrNstda = []*models.NectecStaffTransactionResponse{}
	}

	return &models.GetAllNectecStaffResponse{
		Data:  "All NstdaStaff",
		Obj:   arrNstda,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterNstdaStaff(ctx contractapi.TransactionContextInterface, key, value string) ([]*models.TransactionNectecStaff, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetNstda []*models.TransactionNectecStaff
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset models.TransactionNectecStaff
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetNstda = append(assetNstda, &asset)
		}
	}

	sort.Slice(assetNstda, func(i, j int) bool {
        t1, err1 := time.Parse(time.RFC3339, assetNstda[i].CreatedAt)
        t2, err2 := time.Parse(time.RFC3339, assetNstda[j].CreatedAt)
        if err1 != nil || err2 != nil {
            fmt.Println("Error parsing time:", err1, err2)
            return false
        }
        return t1.After(t2)
    })

	return assetNstda, nil
}