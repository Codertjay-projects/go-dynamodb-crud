package product

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	Validation "gitub.com/go-ozzo/ozzo-validation/v4"
	is "gitub.com/go-ozzo/ozzo-validation/v4/is"
	entities "go-dynamodb-crud/internal/entities"
	product "go-dynamodb-crud/internal/entities/product"
	"io"
	"strings"
	"time"
)

type Rules struct {
}

func NewRules() *Rules {
	return &Rules{}
}

func (r *Rules) ConvertIoReaderToStruct(data io.Reader,
	model interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("body is invalid")
	}
	return model, json.NewDecoder(data).Decode(model)
}
func (r *Rules) Migrate(connection *dynamodb.DynamoDB) error {
	return r.CreateTable(connection)
}

func (r *Rules) GetMock() interface{} {
	return product.Product{
		Base: entities.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: uuid.New().String(),
	}
}

func (r *Rules) Validate(model interface{}) error {
	productModel, err := product.InterfaceToModel(model)
	if err != nil {
		return err
	}
	return Validation.ValidateStruct(productModel,
		Validation.Field(&productModel.ID, Validation.Required, is.UUIDv4),
		Validation.Field(&productModel.Name, Validation.Required),
	)
}

func (r *Rules) CreateTable(connection *dynamodb.DynamoDB) error {
	table := &product.Product{}
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("_id"),
				AttributeType: aws.String("s"),
			},
		},
		BillingMode:            nil,
		GlobalSecondaryIndexes: nil,
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("_id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(table.TableName()),
	}
	response, err := connection.CreateTable(input)
	if err != nil && strings.Contains(err.Error(), "Table already exists") {
		return err
	}
	if response != nil && strings.Contains(
		response.GoString(), "TableStatus:\"CREATING\"") {
		time.Sleep(3 * time.Second)
		err = r.CreateTable(connection)
		if err != nil {
			return err
		}
	}
	return err
}
