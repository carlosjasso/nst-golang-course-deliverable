package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Dynamo struct {
}

func getClient() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("!ERROR: ", err)
	}
	client := dynamodb.NewFromConfig(cfg)
}

func GetTables() ([]string, error) {

}
