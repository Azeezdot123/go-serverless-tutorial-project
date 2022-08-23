package user

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// defining error types I want to send as response
var(
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorFailedToUnMarshalRecord = "failed to unmarshal record"
)

type User struct{
	Email		string	`json:"email"`
	FirstName	string	`json:"firstname"`
	LastName	string	`json:"lastname"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error){
	// query the dynamodb table
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			}
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	// Unmarshal result.Item coming from dynamodb to item of User struct type
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnMarshalRecord)
	}
	return item, nil
}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User, error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	items := new([]User)
	err = dynamodbattribute.UnmarshalMap(result.Items, items)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnMarshalRecord)
	}
	return items, nil
}

func CreateUser()(){}
func UpdateUser()(){}
func DeleteUser() error{}