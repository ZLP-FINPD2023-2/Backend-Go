package domains

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finapp/models"
)

type GoalService interface {
	WithTrx(trxHandle *gorm.DB) GoalService
	List(userID uint) ([]models.GoalCalc, error)
	Create(request *models.GoalCreateRequest, userID uint) error
	Patch(req models.GoalPatchRequest, userID uint) error
	Delete(c *gin.Context, userID uint) error
}
