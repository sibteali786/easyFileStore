package database

import (
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TABLE_NAME = "Users"
)

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}
type UserStore interface {
	DoesUserExists(username string) (bool, error)
	InsertUser(user types.User) error
	GetUser(username string) (types.User, error)
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)
	return DynamoDBClient{
		databaseStore: db,
	}
}

func (u DynamoDBClient) DoesUserExists(username string) (bool, error) {
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return true, err // catch this error in ApiHandler
	}
	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (u DynamoDBClient) InsertUser(user types.User) error {
	// assemble item into a type that dynamoDB understands
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}
	// insert the item
	_, err := u.databaseStore.PutItem(item)
	if err != nil {
		return err
	}
	return nil
}

func (u DynamoDBClient) GetUser(username string) (types.User, error) {
	var user types.User
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}

	// unmarshal the item into our User struct
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, err
	}

	return user, nil

}
