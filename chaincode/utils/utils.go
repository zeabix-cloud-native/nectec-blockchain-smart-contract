package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	UNAUTHORIZE   string = "client is not authorized this asset"
	TIMEFORMAT    string = "2006-01-02T15:04:05Z"
	SKIPOVER      string = "skip over total data"
	DATAUNMARSHAL string = "unmarshal json string"
)

type SmartContract struct {
	contractapi.Contract
}

type GetAllType struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

func Unmarshal(args string, entityType interface{}) (interface{}, error) {
	entityValue := reflect.New(reflect.TypeOf(entityType)).Interface()
	err := json.Unmarshal([]byte(args), entityValue)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json string: %v", err)
	}
	return entityValue, nil
}

func BuildQueryString(filter map[string]interface{}) (string, error) {

	selector := map[string]interface{}{
		"selector": filter,
	}
	queryString, err := json.Marshal(selector)
	if err != nil {
		return "", err
	}
	return string(queryString), nil
}

func CountTotalResults(ctx contractapi.TransactionContextInterface, queryString string) (int, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return 0, err
	}
	defer resultsIterator.Close()

	total := 0
	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return 0, err
		}
		total++
	}
	return total, nil
}

func GetTimeNow() time.Time {
	formattedTime := time.Now().Format(TIMEFORMAT)
	CreatedAt, _ := time.Parse(TIMEFORMAT, formattedTime)
	return CreatedAt
}

func ReturnError(data string) error {
	return fmt.Errorf(data)
}

func GetIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}
func AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func HandleError(err error) error {
	if err != nil {
		return err
	}
	return nil
}
