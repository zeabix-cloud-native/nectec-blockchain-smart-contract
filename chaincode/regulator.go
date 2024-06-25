package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/models"
	"github.com/zeabix-cloud-native/nectec-blockchain-smart-contract/chaincode/utils"
)

type RegulatorContract struct {
    models.SmartContract
}

func (s RegulatorContract) CreateRegulator(
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