package domains

import (
	"gorm.io/gorm"

	"finapp/models"
)

type GoalService interface {
	WithTrx(trxHandle *gorm.DB) GoalService
	List(userID uint) ([]models.GoalCalc, error)
	Create(request *models.GoalCreateRequest, userID uint) error
}
