package domains

import (
	"finapp/models"

	"gorm.io/gorm"
)

type UserService interface {
	WithTrx(trxHandle *gorm.DB) UserService
	GetUserByEmail(email *string) (*models.User, error)
	Register(q *models.RegisterRequest) error
	Delete(id uint) error
	Get(id uint) (*models.User, error)
}
