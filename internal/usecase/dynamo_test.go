package usecase

import (
	"bonded/internal/infra/db"
	"context"
)

type IDynamoDBUsecase interface {
	DynamoDBTest(ctx context.Context) error
}

type DynamoDBUsecase struct {
	dynamo db.IDynamoDB
}

func NewDynamoDBUsecase(dynamo db.IDynamoDB) IDynamoDBUsecase {
	return &DynamoDBUsecase{dynamo: dynamo}
}

func (u *DynamoDBUsecase) DynamoDBTest(ctx context.Context) error {
	return u.dynamo.DescribeSampleTable(ctx)
}
