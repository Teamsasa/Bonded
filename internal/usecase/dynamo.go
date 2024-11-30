package usecase

import (
	"bonded/internal/infra/db"
	"context"
)

type IDynamoUsecase interface {
	DynamoDBTest(ctx context.Context) error
}

type DynamoUsecase struct {
	dynamo db.IDynamoDB
}

func NewDynamoUsecase(dynamo db.IDynamoDB) IDynamoUsecase {
	return &DynamoUsecase{dynamo: dynamo}
}

func (u *DynamoUsecase) DynamoDBTest(ctx context.Context) error {
	return u.dynamo.DescribeSampleTable(ctx)
}
