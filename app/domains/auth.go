package domains

import "finapp/models"

type AuthService interface {
	Authorize(tokenString string) (bool, error)
	CreateToken(user *models.User) (string, error)
	GetTokenClaims(tokenString string) (*models.TokenClaims, error)
}
