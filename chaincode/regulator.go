package chaincode

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

type RegulatorContract struct {
    models.SmartContract
}

func (s RegulatorContract) CreateRegulatorProfile(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
    entityRegulator := models.TransactionRegulator{}
    inputInterface, err := utils.Unmarshal(args, entityRegulator)
	if err != nil {
		return err
	}

	input := inputInterface.(*models.TransactionRegulator)
    orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return utils.ReturnError(utils.UNAUTHORIZE)
	}

    existRegulator, err := utils.AssetExists(ctx, input.Id)
	utils.HandleError(err)
	if existRegulator {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

    clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	CreatedR := utils.GetTimeNow()

	asset := models.TransactionRegulator{
		Id:        input.Id,
		CertId:    input.CertId,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: CreatedR,
		CreatedAt: CreatedR,
	}
	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s RegulatorContract) UpdateRegulatorProfile(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityRegulator := models.TransactionRegulator{}
	inputInterface, err := utils.Unmarshal(args, entityRegulator)
	utils.HandleError(err)
	input := inputInterface.(*models.TransactionRegulator)

	asset, err := s.ReadRegulatorProfile(ctx, input.Id)
	utils.HandleError(err)

	clientID, err := utils.GetIdentity(ctx)
	utils.HandleError(err)

	if clientID != asset.Owner {
		return utils.ReturnError(utils.UNAUTHORIZE)
	}

	UpdatedR := utils.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedR

	assetJSON, err := json.Marshal(asset)
	utils.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *RegulatorContract) ReadRegulatorProfile(ctx contractapi.TransactionContextInterface, id string) (*models.TransactionRegulator, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset models.TransactionRegulator
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *RegulatorContract) GetAllRegulator(ctx contractapi.TransactionContextInterface, args string) (*models.RegulatorGetAllResponse, error) {

	var filterRegulator = map[string]interface{}{}

	entityGetAll := models.FilterGetAllRegulator{}
	interfaceRegulator, err := utils.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := interfaceRegulator.(*models.FilterGetAllRegulator)

	queryStringRegulator, err := utils.BuildQueryString(filterRegulator)
	if err != nil {
		return nil, err
	}

	total, err := utils.CountTotalResults(ctx, queryStringRegulator)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, utils.ReturnError(utils.SKIPOVER)
	}

	arrRegulator, err := utils.RegulatorFetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrRegulator, func(i, j int) bool {
		return arrRegulator[i].UpdatedAt.Before(arrRegulator[j].UpdatedAt)
	})

	if len(arrRegulator) == 0 {
		arrRegulator = []*models.RegulatorTransactionResponse{}
	}

	return &models.RegulatorGetAllResponse{
		Data:  "All Regulator",
		Obj:   arrRegulator,
		Total: total,
	}, nil
}