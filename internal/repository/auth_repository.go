package repository

import (
	"bonded/internal/infra/db"
	"bonded/internal/models"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type IAuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByID(ctx context.Context, userID string) (bool, error)
}

type authRepository struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

func NewAuthRepository(dynamoClient *db.DynamoDBClient) IAuthRepository {
	return &authRepository{
		dynamoDB:  dynamoClient.Client,
		tableName: "Users",
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}
	_, err = r.dynamoDB.PutItemWithContext(ctx, input)
	return err
}

func (r *authRepository) FindUserByID(ctx context.Context, userID string) (bool, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
		},
	}

	result, err := r.dynamoDB.GetItemWithContext(ctx, input)
	if err != nil {
		return false, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}
