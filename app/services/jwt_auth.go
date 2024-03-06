package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"finapp/domains"
	"finapp/lib"
	"finapp/models"
)

// JWTAuthService service relating to authorization
type JWTAuthService struct {
	logger lib.Logger
	env    lib.Env
}

// NewJWTAuthService creates a new auth service
func NewJWTAuthService(logger lib.Logger, env lib.Env) domains.AuthService {
	return JWTAuthService{
		logger: logger,
		env:    env,
	}
}

// Authorize authorizes the generated token
func (s JWTAuthService) Authorize(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.env.SecretKey), nil
	})

	if err != nil {
		return false, err
	}

	if token.Valid {
		return true, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, errors.New("Token malformed")
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, errors.New("Token expired")
		}
	}
	return false, errors.New("Couldn't handle token")
}

// CreateToken creates jwt auth token
func (s JWTAuthService) CreateToken(user *models.User) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + 60*15

	claims := models.TokenClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.env.SecretKey))

	if err != nil {
		s.logger.Error("Failed to sign token string")
		return "", err
	}

	return tokenString, nil
}

func (s JWTAuthService) GetTokenClaims(tokenString string) (*models.TokenClaims, error) {
	claims := &models.TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.env.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*models.TokenClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}
