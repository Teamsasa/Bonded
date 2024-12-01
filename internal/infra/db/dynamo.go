package db

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type IDynamoDB interface {
	DescribeSampleTable(ctx context.Context) error
}

type DynamoDB struct{}

func NewDynamoDB() IDynamoDB {
	return &DynamoDB{}
}

func (db *DynamoDB) DescribeSampleTable(ctx context.Context) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String(os.Getenv("DYNAMODB_ENDPOINT")),
	}))

	svc := dynamodb.New(sess)

	inputCalendars := &dynamodb.DescribeTableInput{
		TableName: aws.String("Calendars"),
	}
	inputEvents := &dynamodb.DescribeTableInput{
		TableName: aws.String("Events"),
	}

	_, err := svc.DescribeTableWithContext(ctx, inputCalendars)
	if err != nil {
		return err
	}
	_, err = svc.DescribeTableWithContext(ctx, inputEvents)
	return err
}
