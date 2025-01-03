package repository

import (
	"bonded/internal/infra/db"
	"bonded/internal/models"
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type eventRepository struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

func EventRepositoryRequest(dynamoClient *db.DynamoDBClient) EventRepository {
	return &eventRepository{
		dynamoDB:  dynamoClient.Client,
		tableName: "Calendars",
	}
}

type calendarRepository struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

func CalendarRepositoryRequest(dynamoClient *db.DynamoDBClient) CalendarRepository {
	return &calendarRepository{
		dynamoDB:  dynamoClient.Client,
		tableName: "Calendars",
	}
}

type CalendarRepository interface {
	Create(ctx context.Context, calendar *models.Calendar) error
	Edit(ctx context.Context, calendarID *models.Calendar, input *models.Calendar) error
	Delete(ctx context.Context, calendarID string) error
	FindAllCalendars(ctx context.Context) ([]*models.Calendar, error)
	FindByCalendarID(ctx context.Context, calendarID string) (*models.Calendar, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.Calendar, error)
	UnfollowCalendar(ctx context.Context, calendar *models.Calendar, user *models.User) error
	FollowCalendar(ctx context.Context, calendar *models.Calendar, user *models.User) error
	InviteUser(ctx context.Context, calendar *models.Calendar, user *models.User) error
}

type EventRepository interface {
	CreateEvent(ctx context.Context, calendar *models.Calendar, event *models.Event) error
	FindEvents(ctx context.Context, calendarID string) ([]*models.Event, error)
	EventExists(ctx context.Context, calendarID string, eventID string) bool
	EditEvent(ctx context.Context, calendarID string, event *models.Event) (*models.Event, error)
	DeleteEvent(ctx context.Context, calendarID string, eventID string) error
}

type userRepository struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

func UserRepositoryRequest(dynamoClient *db.DynamoDBClient) UserRepository {
	return &userRepository{
		dynamoDB:  dynamoClient.Client,
		tableName: "Calendars",
	}
}

type UserRepository interface {
	FindByUserID(ctx context.Context, userID string) (*models.User, error)
}
