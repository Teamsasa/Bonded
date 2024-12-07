#!/bin/bash

TABLE_NAME="Calendars"

echo "Initializing DynamoDB Local..."

# テーブルの存在確認
aws dynamodb describe-table --table-name "$TABLE_NAME" --endpoint-url http://localhost:8000 --region ap-northeast-1 > /dev/null 2>&1

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
        --endpoint-url http://localhost:8000 \
        --region ap-northeast-1 \
        > create_table.log 2>&1

    if [ $? -ne 0 ]; then
        echo "Failed to create table '$TABLE_NAME'. Check create_table.log for details."
        exit 1
    fi

    echo "Waiting for the table to be active..."
    aws dynamodb wait table-exists --table-name "$TABLE_NAME" --endpoint-url http://localhost:8000 --region ap-northeast-1 \
        >> create_table.log 2>&1

    if [ $? -ne 0 ]; then
        echo "Table '$TABLE_NAME' did not become active. Check wait_table.log for details."
        exit 1
    fi

    echo "Table '$TABLE_NAME' created successfully."

    # データ挿入
    echo "Inserting data into '$TABLE_NAME'..."

    # カレンダー情報の挿入
    CALENDARS=("1" "2")
    for CALENDAR in "${CALENDARS[@]}"; do
        aws dynamodb put-item \
            --table-name "$TABLE_NAME" \
            --item "{
                \"CalendarID\": {\"S\": \"$CALENDAR\"},
                \"SortKey\": {\"S\": \"CALENDAR\"},
                \"Name\": {\"S\": \"Test Calendar $CALENDAR\"},
                \"IsPublic\": {\"BOOL\": true},
                \"OwnerUserID\": {\"S\": \"user1\"}
            }" \
            --endpoint-url http://localhost:8000 \
            --region ap-northeast-1 \
            >> insert_data.log 2>&1

        # ユーザー情報の挿入
        if [ "$CALENDAR" == "1" ]; then
            USERS=("user1")
        else
            USERS=("user1" "user2")
        fi

        for USER in "${USERS[@]}"; do
            aws dynamodb put-item \
                --table-name "$TABLE_NAME" \
                --item "{
                    \"CalendarID\": {\"S\": \"$CALENDAR\"},
                    \"SortKey\": {\"S\": \"USER#$USER\"},
                    \"UserID\": {\"S\": \"$USER\"},
                    \"DisplayName\": {\"S\": \"${USER}の表示名\"},
                    \"AccessLevel\": {\"S\": \"OWNER\"}
                }" \
                --endpoint-url http://localhost:8000 \
                --region ap-northeast-1 \
                >> insert_data.log 2>&1

            # GSI用のカレンダー情報の挿入を修正
            aws dynamodb put-item \
                --table-name "$TABLE_NAME" \
                --item "{
                    \"CalendarID\": {\"S\": \"$CALENDAR\"},
                    \"SortKey\": {\"S\": \"CAL#$CALENDAR#$USER\"},
                    \"UserID\": {\"S\": \"$USER\"}
                }" \
                --endpoint-url http://localhost:8000 \
                --region ap-northeast-1 \
                >> insert_data.log 2>&1
        done

        # イベント情報の挿入
        EVENTS=("event1" "event2" "event3")
        for EVENT in "${EVENTS[@]}"; do
            aws dynamodb put-item \
                --table-name "$TABLE_NAME" \
                --item "{
                    \"CalendarID\": {\"S\": \"$CALENDAR\"},
                    \"SortKey\": {\"S\": \"EVENT#$EVENT\"},
                    \"EventID\": {\"S\": \"$EVENT\"},
                    \"Title\": {\"S\": \"Test Event $EVENT\"},
                    \"Description\": {\"S\": \"This is test event $EVENT\"},
                    \"StartTime\": {\"S\": \"2021-08-0${EVENT: -1}T00:00:00Z\"},
                    \"EndTime\": {\"S\": \"2021-08-0${EVENT: -1}T01:00:00Z\"},
                    \"Location\": {\"S\": \"場所$EVENT\"},
                    \"AllDay\": {\"BOOL\": false}
                }" \
                --endpoint-url http://localhost:8000 \
                --region ap-northeast-1 \
                >> insert_data.log 2>&1
        done
    done

    echo "Data inserted successfully."

else
    echo "Table '$TABLE_NAME' already exists. Skipping table creation."
fi

echo "DynamoDB Initialization Complete."
