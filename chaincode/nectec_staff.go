package chaincode

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

func (s *SmartContract) CreateNstdaStaff(
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

	TimeNstda := utils.GetTimeNow()

	asset := models.TransactionNectecStaff{
		Id:        input.Id,
		CertId:    input.CertId,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: TimeNstda,
		CreatedAt: TimeNstda,
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

	UpdatedNstda := utils.GetTimeNow()
	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedNstda

	assetJSON, errN := json.Marshal(asset)
	utils.HandleError(errN)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteNectecStaff(ctx contractapi.TransactionContextInterface, id string) error {

	assetNstda, err := s.ReadNectecStaff(ctx, id)
	utils.HandleError(err)

	clientIDNstda, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientIDNstda != assetNstda.Owner {
		return fmt.Errorf(utils.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
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

func (s *SmartContract) GetAllNstdaStaff(ctx contractapi.TransactionContextInterface, args string) (*models.GetAllNectecStaffResponse, error) {

	var filterNstda = map[string]interface{}{}

	entityGetAll := models.FilterGetAllNectecStaff{}
	interfaceNstda, err := utils.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := interfaceNstda.(*models.FilterGetAllNectecStaff)

	queryStringNstda, err := utils.BuildQueryString(filterNstda)
	if err != nil {
		return nil, err
	}

	total, err := utils.CountTotalResults(ctx, queryStringNstda)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, utils.ReturnError(utils.SKIPOVER)
	}

	arrNstda, err := utils.NectecStaffFetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrNstda, func(i, j int) bool {
		return arrNstda[i].UpdatedAt.Before(arrNstda[j].UpdatedAt)
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
		return assetNstda[i].UpdatedAt.After(assetNstda[j].UpdatedAt)
	})

	return assetNstda, nil
}