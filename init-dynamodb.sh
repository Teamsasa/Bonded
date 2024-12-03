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
            AttributeName=UserID,AttributeType=S \
            AttributeName=CalendarID,AttributeType=S \
        --key-schema \
            AttributeName=UserID,KeyType=HASH \
            AttributeName=CalendarID,KeyType=RANGE \
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
      --item '{
          "UserID": {"S": "user1"},
          "CalendarID": {"S": "1"},
          "Name": {"S": "Test User"},
          "Event": {
              "L": [
                  {
                      "M": {
                          "EventID": {"S": "1"},
                          "Title": {"S": "Test Event 1"},
                          "Description": {"S": "This is a test event 1"},
                          "StartTime": {"S": "2021-08-01T00:00:00Z"},
                          "EndTime": {"S": "2021-08-01T01:00:00Z"},
                          "Location": {"S": "新宿"},
                          "AllDay": {"BOOL": false}
                      }
                  },
                  {
                      "M": {
                          "EventID": {"S": "2"},
                          "Title": {"S": "Test Event 2"},
                          "Description": {"S": "This is a test event 2"},
                          "StartTime": {"S": "2021-08-02T00:00:00Z"},
                          "EndTime": {"S": "2021-08-02T01:00:00Z"},
                          "Location": {"S": "渋谷"},
                          "AllDay": {"BOOL": false}
                      }
                  }
              ]
          }
      }' \
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

echo "DynamoDB Initialization Complete."