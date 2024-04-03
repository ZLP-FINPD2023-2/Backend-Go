package domains

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finapp/models"
)

type BudgetService interface {
	WithTrx(trxHandle *gorm.DB) BudgetService
	List(userID uint) ([]models.BudgetCalc, error)
	Create(request *models.BudgetCreateRequest, userID uint) error
	Patch(budget models.BudgetPatchRequest, userID uint) error
	Delete(c *gin.Context, userID uint) error
}
