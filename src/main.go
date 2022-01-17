package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type AwsSettings struct {
	Key    string `json:"aws_access_key_id"`
	Secret string `json:"aws_secret_access_key"`
}

func main() {
	jsonSettingsFile, err := os.Open("go.secrets.json")
	if err != nil {
		log.Fatalf("unable to open config file: %v", err)
	}

	jsonSettingsBytes, err := ioutil.ReadAll(jsonSettingsFile)
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}

	var awsSettings AwsSettings
	err = json.Unmarshal(jsonSettingsBytes, &awsSettings)
	if err != nil {
		log.Fatalf("unable to deserialize config file: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(awsSettings.Key, awsSettings.Secret, ""),
		))
	if err != nil {
		log.Fatalf("unable to load aws configuration: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	resp, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("could not list dynamodb tables: %v", err)
	}

	fmt.Println("Tables: ", resp.TableNames)

	tableName := "<table>"
	uuid := uuid.New().String()

	output, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: uuid,
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatalf("could not put item: %v", err)
	}

	fmt.Println("dynamo insert output: ", output)

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": types.AttributeValueMemberS{
				Value: uuid,
			},
		},
	})
	if err != nil {
		log.Fatalf("could not read item: %v", err)
	}

	if result.Item == nil {
		fmt.Println("Item with id %v not found", uuid)
	} else {
		fmt.Println("Item with id %v: ", result.Item)
	}
}
