package usecase

import (
	"errors"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

type IAuthUsecase interface {
	ValidateJWT(tokenString string) (*jwt.Token, error)
}

type AuthUsecase struct { 
	jwks          *keyfunc.JWKS
	clientID      string
	cognitoIssuer string
	nameGenSeed   int64
}

func NewAuthUsecase(
	jwks *keyfunc.JWKS,
	clientID string,
	cognitoIssuer string,
) *AuthUsecase {
	return &AuthUsecase{
		jwks:          jwks,
		clientID:      clientID,
		cognitoIssuer: cognitoIssuer,
		nameGenSeed:   time.Now().UTC().UnixNano(),
	}
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
