#!/bin/bash

TABLE_NAME="Calendars"
ENDPOINT_URL="http://localhost:8000"
REGION="us-west-2"

echo "Initializing DynamoDB Local..."

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
        --endpoint-url "$ENDPOINT_URL" \
        --region "$REGION" \
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
                \"Email\": {\"S\": \"${USER}@example.com\"},
                \"Password\": {\"S\": \"password\"},
                \"AccessLevel\": {\"S\": \"OWNER\"}
            }" \
            --endpoint-url "$ENDPOINT_URL" \
            --region "$REGION" \
            >> insert_data.log 2>&1

        # GSI用のカレンダー情報の挿入を修正
        aws dynamodb put-item \
            --table-name "$TABLE_NAME" \
            --item "{
                \"CalendarID\": {\"S\": \"$CALENDAR\"},
                \"SortKey\": {\"S\": \"CAL#$CALENDAR#$USER\"},
                \"UserID\": {\"S\": \"$USER\"}
            }" \
            --endpoint-url "$ENDPOINT_URL" \
            --region "$REGION" \
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
            --endpoint-url "$ENDPOINT_URL" \
            --region "$REGION" \
            >> insert_data.log 2>&1
    done
done
    echo "Data inserted successfully."
fi

echo "DynamoDB Initialization Complete."
