package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/azeezdot123/go-serverless/pkg/validators"
)

// defining error types I want to send as response
var(
	ErrorFailedToFetchRecord		= "failed to fetch record"
	ErrorFailedToUnMarshalRecord	= "failed to unmarshal record"
	ErrorInvalidUserData			= "Invalid user data"
	ErrorInvalidEmail				= " Invalid email"
	ErrorCouldNotMarshalItem 		= "failed to marshal item"
	ErrorCouldNotDeleteItem 		= "could not delete"
	ErrorCouldNotDynamoPutItem 		= "could not dynamo put item"
	ErrorUserAlreadyExists			= "user already exists"
	ErrorUserDoesNotExist			= "user does not exist"
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

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error){
	var u User

	if err := json.Unmarshal([]byte(req.body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}
	if !validators.IsEmailValid(u.Email){
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	av, err := dynamodbattribute.marshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil
}

func UpdateUser()(){}
func DeleteUser() error{}