package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBClient struct {
	Client *dynamodb.DynamoDB
}

func DynamoDBClientRequest() *DynamoDBClient {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000" // デフォルトエンドポイント
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String(endpoint),
	}))
	return &DynamoDBClient{
		Client: dynamodb.New(sess),
	}
}
