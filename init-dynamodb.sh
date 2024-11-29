#!/bin/bash

TABLE_NAME="sampleTable"

echo "Initializing DynamoDB Local..."

# テーブルの存在確認
if aws dynamodb list-tables \
  --endpoint-url http://localhost:8000 \
  --region ap-northeast-1 \
  --query "TableNames" \
  --output text | grep -w "$TABLE_NAME" > /dev/null 2>&1; then
  echo "Table '$TABLE_NAME' already exists. Skipping table creation."
else
  echo "Creating table '$TABLE_NAME'..."
  aws dynamodb create-table \
    --table-name "$TABLE_NAME" \
    --attribute-definitions AttributeName=Id,AttributeType=S \
    --key-schema AttributeName=Id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:8000 \
    --region ap-northeast-1 > /dev/null 2>&1
  echo "Table '$TABLE_NAME' created."
fi

# データ挿入
echo "Inserting data..."
aws dynamodb put-item \
  --table-name "$TABLE_NAME" \
  --item '{"Id": {"S": "1"}, "Name": {"S": "testUser1"}}' \
  --endpoint-url http://localhost:8000 \
  --region ap-northeast-1 > /dev/null 2>&1
echo "Data inserted successfully."

echo "DynamoDB Initialization Complete."

