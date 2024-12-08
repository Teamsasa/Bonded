package repository

import (
	"bonded/internal/models"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (r *userRepository) FindByUserID(ctx context.Context, userID string) (*models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("UserID-index"),
		KeyConditionExpression: aws.String("UserID = :uid"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":uid": {S: aws.String(userID)},
		},
	}

	result, err := r.dynamoDB.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user with UserID %s not found", userID)
	}

	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.DisplayName != "" && user.AccessLevel != "" {
			fmt.Printf("Selected user: %v\n", user)
			return &user, nil
		}
	}

	return nil, fmt.Errorf("valid user with UserID %s not found", userID)
}
