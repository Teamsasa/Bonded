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
            AttributeName=CalendarID,KeyType=HASH \
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
          "CalendarID": {"S": "1"},
          "Name": {"S": "Test Calendar 1"},
          "IsPublic": {"BOOL": true},
          "OwnerUserID": {"S": "user1"},
          "Users": {
              "L": [
                  {
                      "M": {
                          "UserID": {"S": "user1"},
                          "DisplayName": {"S": "ユーザー1の表示名"},
                          "Email": {"S": "user1@example.com"},
                          "Password": {"S": "password1"},
                          "AccessLevel": {"S": "OWNER"}
                      }
                  },
                  {
                      "M": {
                          "UserID": {"S": "user2"},
                          "DisplayName": {"S": "ユーザー2の表示名"},
                          "Email": {"S": "user2@example.com"},
                          "Password": {"S": "password2"},
                          "AccessLevel": {"S": "EDITOR"}
                      }
                  }
              ]
          },
          "Events": {
              "L": [
                  {
                      "M": {
                          "EventID": {"S": "event1"},
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
                          "EventID": {"S": "event2"},
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