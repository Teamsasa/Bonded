package usecase

import (
	"bonded/internal/infra/db"
	"context"
)

func DynamoUsecaseRequest(dynamo db.IDynamoDB) IDynamoUsecase {
	return &DynamoUsecase{dynamo: dynamo}
}

type IDynamoUsecase interface {
	DynamoDBTest(ctx context.Context) error
}

type DynamoUsecase struct {
	dynamo db.IDynamoDB
}

func (u *DynamoUsecase) DynamoDBTest(ctx context.Context) error {
	return u.dynamo.DescribeSampleTable(ctx)
}
