package aws

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoItem struct {
	Id string
}

func client(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

func listTables(cfg aws.Config) []string {
	resp, err := client(cfg).ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("could not list dynamodb tables: %v", err)
	}
	return resp.TableNames
}

func getDefaultTableName(cfg aws.Config) *string {
	var tableName string
	for _, t := range listTables(cfg) {
		if t == "MyTable" {
			tableName = t
			break
		}
	}
	if tableName == "" {
		log.Fatalln(errors.New("table with name \"MyTable\" not found"))
	}
	return aws.String(tableName)
}

func Insert(cfg aws.Config, value string) bool {
	_, err := client(cfg).PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{
				Value: value,
			},
		},
		TableName: getDefaultTableName(cfg),
	})
	if err != nil {
		log.Println("!ERROR could not put item - ", err)
		return false
	}

	return true
}

/*
AWS Dynamodb read doc:
*/

func Read(cfg aws.Config, value string) bool {
	getItemOutput, err := client(cfg).GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: getDefaultTableName(cfg),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{
				Value: value,
			},
		},
	})
	if err != nil {
		log.Printf("!ERROR could not read item: %v. %v\n", value, err)
		return false
	}
	if getItemOutput.Item == nil {
		log.Printf("!WARNING item with id %v not found\n", value)
		return false
	} else {
		var item DynamoItem
		err := attributevalue.UnmarshalMap(getItemOutput.Item, &item)
		if err != nil {
			log.Printf("!ERROR could not unmarshal item: %v\n", err)
			return false
		}
		return item.Id == value
	}
}
