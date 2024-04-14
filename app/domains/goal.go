package domains

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finapp/models"
)

type GoalService interface {
	WithTrx(trxHandle *gorm.DB) GoalService
	List(userID uint) ([]models.GoalCalc, error)
	Store(request *models.GoalStoreRequest, userID uint) error
	Update(req models.GoalUpdateRequest, userID uint) error
	Delete(c *gin.Context, userID uint) error
}
