package adapter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Database struct {
	connection *dynamodb.DynamoDB
	logMode    bool
}

// Interface this interface contains all the functions used to
//connect,filter, delete and more through the db
type Interface interface {
	Health() bool
	FindAll(condition expression.Expression,
		tableName string) (response *dynamodb.ScanOutput,
		err error)
	FindOne(condition map[string]interface{}, tableName string) (
		response *dynamodb.GetItemOutput, err error)
	CreateOrUpdate(entity interface{}, tableName string) (
		response *dynamodb.PutItemOutput, err error)
	Delete(condition map[string]interface{}, tableName string) (
		response *dynamodb.DeleteItemOutput, err error)
}

// NewAdapter database takes in a connection gotten from
//GetConnection in our instance
func NewAdapter(con *dynamodb.DynamoDB) Interface {
	return &Database{
		connection: con,
		logMode:    false,
	}
}

func (db *Database) Health() bool {
	_, err := db.connection.ListTables(&dynamodb.ListTablesInput{})
	return err == nil
}
func (db *Database) FindAll(
	condition expression.Expression,
	tableName string) (response *dynamodb.ScanOutput, err error) {
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: condition.Names(),
		ExclusiveStartKey:        condition.Values(),
		FilterExpression:         condition.Filter(),
		ProjectionExpression:     condition.Projection(),
		TableName:                aws.String(tableName),
	}
	return db.connection.Scan(input)
}

func (db *Database) FindOne(condition map[string]interface{},
	tableName string) (response *dynamodb.GetItemOutput, err error) {

	conditionParsed, err := dynamodbattribute.MarshalMap(condition)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       conditionParsed,
	}
	return db.connection.GetItem(input)
}

func (db *Database) CreateOrUpdate(
	entity interface{},
	tableName string) (response *dynamodb.PutItemOutput, err error) {

	// converting the struct product to map for dynamodb to understand
	entityParsed, err := dynamodbattribute.MarshalMap(entity)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(tableName),
	}
	return db.connection.PutItem(input)
}

func (db *Database) Delete(
	condition map[string]interface{},
	tableName string) (response *dynamodb.DeleteItemOutput, err error) {

	conditionParsed, err := dynamodbattribute.MarshalMap(condition)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.DeleteItemInput{
		Key:       conditionParsed,
		TableName: aws.String(tableName),
	}
	deleteItemOutput, err := db.connection.DeleteItem(input)
	if err != nil {
		return nil, err
	}
	return deleteItemOutput, nil
}
