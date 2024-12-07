#!/bin/bash

TABLE_NAME="Calendars"
ENDPOINT_URL="http://dynamodb.us-west-2.amazonaws.com"
REGION="us-west-2"

echo "Initializing DynamoDB..."

# テーブルの存在確認
aws dynamodb describe-table --table-name "$TABLE_NAME" --endpoint-url "$ENDPOINT_URL" --region "$REGION" > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo "Creating DynamoDB table '$TABLE_NAME'..."
    aws dynamodb create-table \
        --table-name "$TABLE_NAME" \
        --attribute-definitions \
            AttributeName=CalendarID,AttributeType=S \
            AttributeName=SortKey,AttributeType=S \
            AttributeName=UserID,AttributeType=S \
        --key-schema \
            AttributeName=CalendarID,KeyType=HASH \
            AttributeName=SortKey,KeyType=RANGE \
        --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
        --global-secondary-indexes \
            "[
                {
                    \"IndexName\": \"UserID-index\",
                    \"KeySchema\": [
                        {\"AttributeName\":\"UserID\",\"KeyType\":\"HASH\"}
                    ],
                    \"Projection\":{
                        \"ProjectionType\":\"ALL\"
                    },
                    \"ProvisionedThroughput\": {
                        \"ReadCapacityUnits\": 5,
                        \"WriteCapacityUnits\": 5
                    }
                }
            ]" \
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" \
        > create_table.log 2>&1

    if [ $? -ne 0 ]; then
        echo "Failed to create table '$TABLE_NAME'. Check create_table.log for details."
        exit 1
    fi

    echo "Waiting for the table to be active..."
    aws dynamodb wait table-exists --table-name "$TABLE_NAME" --endpoint-url "$ENDPOINT_URL" --region "$REGION" \
        >> create_table.log 2>&1

    if [ $? -ne 0 ]; then
        echo "Table '$TABLE_NAME' did not become active. Check wait_table.log for details."
        exit 1
    fi

    echo "Remote Table '$TABLE_NAME' created successfully."
else
    echo "Remote Table '$TABLE_NAME' already exists. Skipping creation."
fi
echo "Remote DynamoDB Initialization Complete."
