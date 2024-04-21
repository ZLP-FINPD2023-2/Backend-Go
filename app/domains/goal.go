package domains

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finapp/models"
)

type GoalService interface {
	WithTrx(trxHandle *gorm.DB) GoalService
	List(c *gin.Context, userID uint) ([]models.GoalCalcResponse, error)
	Get(c *gin.Context, userID uint) (models.GoalCalcResponse, error)
	Store(request *models.GoalStoreRequest, userID uint) (models.GoalResponse, error)
	Update(c *gin.Context, req models.GoalUpdateRequest, userID uint) (models.GoalResponse, error)
	Delete(c *gin.Context, userID uint) error
}
