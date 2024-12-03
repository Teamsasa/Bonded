package usecase

import (
	"bonded/internal/models"
	"bonded/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/goombaio/namegenerator"
)

type IAuthUsecase interface {
	FindOrCreateUserByID(ctx context.Context, userID string) error
	ValidateJWT(tokenString string) (*jwt.Token, error)
}

type AuthUsecase struct {
	authRepo      repository.IAuthRepository
	jwks          *keyfunc.JWKS
	clientID      string
	cognitoIssuer string
	nameGenSeed   int64
}

func NewAuthUsecase(
	authRepo repository.IAuthRepository,
	jwks *keyfunc.JWKS,
	clientID string,
	cognitoIssuer string,
) *AuthUsecase {
	return &AuthUsecase{
		authRepo:      authRepo,
		jwks:          jwks,
		clientID:      clientID,
		cognitoIssuer: cognitoIssuer,
		nameGenSeed:   time.Now().UTC().UnixNano(),
	}
}

func (u *AuthUsecase) FindOrCreateUserByID(ctx context.Context, userID string) error {
	found, err := u.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if found {
		return nil
	}
	nameGenerator := namegenerator.NewNameGenerator(u.nameGenSeed)
	name := nameGenerator.Generate()
	user := &models.User{
		UserId: userID,
		Name:   name,
	}
	return u.authRepo.CreateUser(ctx, user)
}

func (u *AuthUsecase) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, u.jwks.Keyfunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("token expired")
	}

	if !claims.VerifyIssuer(u.cognitoIssuer, true) {
		return nil, errors.New("invalid issuer")
	}

	if !claims.VerifyAudience(u.clientID, false) {
		return nil, errors.New("invalid audience")
	}

	return token, nil
}
