package usecase

import (
	"bonded/internal/contextKey"
	"bonded/internal/models"
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func (u *calendarUsecase) FindCalendar(ctx context.Context, calendarID string) (*models.Calendar, error) {
	calendarData, err := u.calendarRepo.FindByCalendarID(ctx, calendarID)
	if err != nil {
		return nil, err
	}

	if !*calendarData.IsPublic {
		jwtData, ok := ctx.Value(contextKey.JwtDataKey).(*jwt.Token)
		if !ok {
			return nil, errors.New("failed to get JWT data from context")
		}

		accessUserID, ok := jwtData.Claims.(jwt.MapClaims)["sub"].(string)
		if !ok {
			return nil, errors.New("failed to get UserID from JWT data")
		}

		isExist := false
		for _, user := range calendarData.Users {
			if user.UserID == accessUserID {
				isExist = true
				break
			}
		}

		if !isExist {
			return nil, errors.New("access user is not registered in the calendar")
		}
	}

	return calendarData, nil
}

func (u *calendarUsecase) FindPublicCalendars(ctx context.Context) ([]*models.Calendar, error) {
	//全件取得してフィルタリング
	calendars, err := u.calendarRepo.FindAllCalendars(ctx)
	if err != nil {
		return nil, err
	}
	publicCalendars := []*models.Calendar{}
	for _, calendar := range calendars {
		if *calendar.IsPublic {
			publicCalendars = append(publicCalendars, calendar)
		}
	}
	return publicCalendars, nil
}

func (u *calendarUsecase) CreateCalendar(ctx context.Context, calendar *models.CreateCalendar) error {

	jwtData, ok := ctx.Value(contextKey.JwtDataKey).(*jwt.Token)
	if !ok {
		return errors.New("failed to get JWT data from context")
	}

	accessUserID, ok := jwtData.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		return errors.New("failed to get UserID from JWT data")
	}
	calendar.OwnerUserID = accessUserID
	if calendar.OwnerName == "" {
		user, err := u.userRepo.FindByUserID(ctx, calendar.OwnerUserID)
		if err != nil {
			return err
		}
		calendar.OwnerName = user.DisplayName
	}
	user := models.User{
		UserID:      calendar.OwnerUserID,
		DisplayName: calendar.OwnerName,
		Email:       calendar.OwnerUserID + "@example.com",
		Password:    "password",
		AccessLevel: "OWNER",
	}
	calendar.Users = []models.User{user}

	calendar.CalendarID = uuid.New().String()

	// CreateCalendarのフィールドをCalendarに変換
	calendarReq := models.Calendar{
		CalendarID:  calendar.CalendarID,
		SortKey:     "CALENDAR",
		Name:        calendar.Name,
		IsPublic:    calendar.IsPublic,
		OwnerUserID: calendar.OwnerUserID,
		Users:       calendar.Users,
		Events:      calendar.Events,
	}

	return u.calendarRepo.Create(ctx, &calendarReq)
}

func (u *calendarUsecase) EditCalendar(ctx context.Context, calendar *models.Calendar, input *models.Calendar) error {
	return u.calendarRepo.Edit(ctx, calendar, input)
}

func (u *calendarUsecase) DeleteCalendar(ctx context.Context, calendarID string) error {
	return u.calendarRepo.Delete(ctx, calendarID)
}

func (u *calendarUsecase) FindCalendars(ctx context.Context) ([]*models.Calendar, error) {
	jwtData, ok := ctx.Value(contextKey.JwtDataKey).(*jwt.Token)
	if !ok {
		return nil, errors.New("failed to get JWT data from context")
	}

	accessUserID, ok := jwtData.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		return nil, errors.New("failed to get UserID from JWT data")
	}
	return u.calendarRepo.FindByUserID(ctx, accessUserID)
}

func (u *calendarUsecase) FollowCalendar(ctx context.Context, calendar *models.Calendar) error {
	jwtData, ok := ctx.Value(contextKey.JwtDataKey).(*jwt.Token)
	if !ok {
		return errors.New("failed to get JWT data from context")
	}

	accessUserID, ok := jwtData.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		return errors.New("failed to get UserID from JWT data")
	}
	user, err := u.userRepo.FindByUserID(ctx, accessUserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return u.calendarRepo.FollowCalendar(ctx, calendar, user)
}

func (u *calendarUsecase) UnfollowCalendar(ctx context.Context, calendar *models.Calendar) error {
	jwtData, ok := ctx.Value(contextKey.JwtDataKey).(*jwt.Token)
	if !ok {
		return errors.New("failed to get JWT data from context")
	}

	accessUserID, ok := jwtData.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		return errors.New("failed to get UserID from JWT data")
	}
	user, err := u.userRepo.FindByUserID(ctx, accessUserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return u.calendarRepo.UnfollowCalendar(ctx, calendar, user)
}

func (u *calendarUsecase) InviteUser(ctx context.Context, calendarID string, inviteUserID string, accessLevel string) error {
	// カレンダーの取得

	jwtData, ok := ctx.Value(contextKey.JwtDataKey).(*jwt.Token)
	if !ok {
		return errors.New("failed to get JWT data from context")
	}

	ownerUserID, ok := jwtData.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		return errors.New("failed to get UserID from JWT data")
	}

	calendar, err := u.calendarRepo.FindByCalendarID(ctx, calendarID)
	if err != nil {
		return err
	}
	if calendar == nil {
		return errors.New("calendar not found")
	}

	// 非公開カレンダーの場合、オーナーチェック
	if !*calendar.IsPublic && calendar.OwnerUserID != ownerUserID {
		return errors.New("only the owner can invite users to private calendars")
	}

	// 招待するユーザーの存在確認
	inviteUser, err := u.userRepo.FindByUserID(ctx, inviteUserID)
	if err != nil {
		return err
	}
	if inviteUser == nil {
		return errors.New("invite user not found")
	}

	// ユーザーが既に追加されているか確認
	for _, user := range calendar.Users {
		if user.UserID == inviteUserID {
			return errors.New("user is already a member of this calendar")
		}
	}

	// ユーザー情報の設定
	inviteUser.AccessLevel = accessLevel

	return u.calendarRepo.InviteUser(ctx, calendar, inviteUser)
}
