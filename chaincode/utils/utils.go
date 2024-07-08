package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	COUCHDB_DATEFORMAT string = "02-01-2006"
	DATEFORMAT    string = "02-01-2006"
	UNAUTHORIZE   string = "client is not authorized this asset"
	TIMEFORMAT    string = "2006-01-02T15:04:05Z"
	SKIPOVER      string = "skip over total data"
	DATAUNMARSHAL string = "unmarshal json string"
)

const offset = 7 
type SmartContract struct {
	contractapi.Contract
}

type GetAllType struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

func FormatDate(dateStr string, isEndDate bool, offset int) (string, error) {
	const inputFormat = "02-01-2006"
	const outputFormat = time.RFC3339

	// Parse the date in the specified timezone
	location := time.FixedZone("UTC+7", offset*3600)
	parsedDate, err := time.ParseInLocation(inputFormat, dateStr, location)
	if err != nil {
		return "", err
	}

	if isEndDate {
		// Set the time to the end of the day
		parsedDate = parsedDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Convert to UTC
	utcDate := parsedDate.UTC()

	return utcDate.Format(outputFormat), nil
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

func ParseDate(input string) (string, error) {
	sanitizedInput := strings.ReplaceAll(input, "–", "-")

	parsedTime, err := time.Parse(DATEFORMAT, sanitizedInput)
	if err != nil {
		return "", err
	}

	return parsedTime.Format(COUCHDB_DATEFORMAT), nil
}

func SanitizeDate(dateStr string) (string, error) {
	dateStr = strings.ReplaceAll(dateStr, "–", "-")
	parsedDate, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02"), nil
}


func GetTimeNow() time.Time {
	now := time.Now()
	
	formattedTime := now.Format(TIMEFORMAT)
	
	parsedTime, err := time.ParseInLocation(TIMEFORMAT, formattedTime, time.Local)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}
	
	return parsedTime
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

func StructToMap(v interface{}) (map[string]interface{}, error) {
    var m map[string]interface{}
    data, err := json.Marshal(v)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &m)
    if err != nil {
        return nil, err
    }
    return m, nil
}