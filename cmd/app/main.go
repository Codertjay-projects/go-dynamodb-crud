package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go-dynamodb-crud/config"
	"go-dynamodb-crud/internal/repository/adapter"
	"go-dynamodb-crud/internal/repository/instance"
	"go-dynamodb-crud/internal/routes"
	"go-dynamodb-crud/internal/rules"
	RulesProduct "go-dynamodb-crud/internal/rules/product"
	"go-dynamodb-crud/utils/logger"
	"log"
	"net/http"
)

func main() {
	// config i created to set api config
	configs := config.GetConfig()
	// connect to dynamo db
	connection := instance.GetConnection()
	// this basically get the entire database could be used to access it
	repository := adapter.NewAdapter(connection)
	// logger we created
	logger.INFO("waiting for the service to start...", nil)
	errors := Migrate(connection)
	if len(errors) > 0 {
		for _, err := range errors {
			logger.PANIC("Error on migration....", err)
		}
	}
	logger.PANIC("", CheckTables(connection))
	port := fmt.Sprintf(":v%", configs.Port)
	// creating a new router and sending the repository -> router -> handler
	router := routes.NewRouter().SetRouter(repository)
	// list to server
	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.DynamoDB) []error {
	var errors []error
	callMigrateAndAppendError(&errors, connection, &RulesProduct.Rules{})
	return errors
}

func callMigrateAndAppendError(
	errors *[]error,
	connection *dynamodb.DynamoDB, rule rules.Interface) []error {

	err := rule.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
	return *errors
}
func CheckTables(connection *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTablesInput{})
	if len(response.TableNames) == 0 {
		logger.INFO("Tables not found:", nil)
	}
	for _, tableName := range response.TableNames {
		logger.INFO("Table found:", *tableName)
	}
	return err
}
