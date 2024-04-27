package domains

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finapp/models"
)

type BudgetService interface {
	WithTrx(trxHandle *gorm.DB) BudgetService
	List(c *gin.Context, userID uint) ([]models.BudgetGetResponse, error)
	Get(c *gin.Context, userID uint) (models.BudgetGetResponse, error)
	Create(request *models.BudgetCreateRequest, userID uint) (models.BudgetCreateResponse, error)
	Patch(c *gin.Context, budget models.BudgetPatchRequest, userID uint) (models.BudgetPatchResponse, error)
	Delete(c *gin.Context, userID uint) error
}
