#!/bin/bash

TABLE_NAME="Calendars"

echo "Initializing DynamoDB Local..."

# テーブルの存在確認
aws dynamodb describe-table --table-name "$TABLE_NAME" --endpoint-url http://localhost:8000 --region ap-northeast-1 2>/dev/null

if [ $? -ne 0 ]; then
    echo "Creating DynamoDB table '$TABLE_NAME'..."
    aws dynamodb create-table \
        --table-name "$TABLE_NAME" \
        --attribute-definitions \
            AttributeName=CalendarID,AttributeType=S \
            AttributeName=UserID,AttributeType=S \
        --key-schema \
            AttributeName=CalendarID,KeyType=HASH \
        --global-secondary-indexes file://gsi_calendars.json \
        --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
        --endpoint-url http://localhost:8000 \
        --region ap-northeast-1 \
        > create_table.log 2>&1

    if [ $? -ne 0 ]; then
        echo "Failed to create table '$TABLE_NAME'. Check create_table.log for details."
        exit 1
    fi

    echo "Waiting for the table to be active..."
    aws dynamodb wait table-exists --table-name "$TABLE_NAME" --endpoint-url http://localhost:8000 --region ap-northeast-1 \
        > wait_table.log 2>&1

    if [ $? -ne 0 ]; then
        echo "Table '$TABLE_NAME' did not become active. Check wait_table.log for details."
        exit 1
    fi

    echo "Table '$TABLE_NAME' created successfully."

    # データ挿入
    echo "Inserting data into '$TABLE_NAME'..."
    aws dynamodb put-item \
      --table-name "$TABLE_NAME" \
      --item '{"CalendarID": {"S": "1"}, "UserID": {"S": "user1"}, "Name": {"S": "Test User"}}' \
      --endpoint-url http://localhost:8000 \
      --region ap-northeast-1 \
      > insert_data.log 2>&1

    if [ $? -eq 0 ]; then
        echo "Data inserted successfully."
    else
        echo "Failed to insert data. Check insert_data.log for details."
    fi

else
    echo "Table '$TABLE_NAME' already exists. Skipping table creation."
fi

EVENT_TABLE_NAME="Events"

echo "Initializing Events DynamoDB table..."

# テーブルの存在確認
aws dynamodb describe-table --table-name "$EVENT_TABLE_NAME" --endpoint-url http://localhost:8000 --region ap-northeast-1 2>/dev/null

if [ $? -ne 0 ]; then
    echo "Creating DynamoDB table '$EVENT_TABLE_NAME'..."
    aws dynamodb create-table \
        --table-name "$EVENT_TABLE_NAME" \
        --attribute-definitions \
            AttributeName=ID,AttributeType=S \
            AttributeName=CalendarID,AttributeType=S \
        --key-schema \
            AttributeName=ID,KeyType=HASH \
        --global-secondary-indexes file://gsi_events.json \
        --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
        --endpoint-url http://localhost:8000 \
        --region ap-northeast-1 \
        > create_events_table.log 2>&1

    if [ $? -ne 0 ]; then
        echo "Failed to create table '$EVENT_TABLE_NAME'. Check create_events_table.log for details."
        exit 1
    fi

    echo "Waiting for the table to be active..."
    aws dynamodb wait table-exists --table-name "$EVENT_TABLE_NAME" --endpoint-url http://localhost:8000 --region ap-northeast-1 \
        > wait_events_table.log 2>&1

    if [ $? -ne 0 ]; then
        echo "Table '$EVENT_TABLE_NAME' did not become active. Check wait_events_table.log for details."
        exit 1
    fi

    echo "Table '$EVENT_TABLE_NAME' created successfully."

    # データ挿入
    echo "Inserting data into '$EVENT_TABLE_NAME'..."
    aws dynamodb put-item \
      --table-name "$EVENT_TABLE_NAME" \
      --item '{
          "ID": {"S": "1"},
          "CalendarID": {"S": "64f1b9cb-fb79-403a-a4d8-e7203a80f52a"},
          "Title": {"S": "Test Event"},
          "Description": {"S": "This is a test event"},
          "StartTime": {"S": "2021-08-01T00:00:00Z"},
          "EndTime": {"S": "2021-08-01T01:00:00Z"},
          "Location": {"S": "新宿"},
          "AllDay": {"BOOL": false}
      }' \
      --endpoint-url http://localhost:8000 \
      --region ap-northeast-1 \
      > insert_event_data.log 2>&1

    if [ $? -eq 0 ]; then
        echo "Data inserted successfully."
    else
        echo "Failed to insert data. Check insert_event_data.log for details."
    fi

else
    echo "Table '$EVENT_TABLE_NAME' already exists. Skipping table creation."
fi

echo "DynamoDB Initialization Complete."