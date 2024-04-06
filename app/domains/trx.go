package domains

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finapp/models"
)

type TrxService interface {
	WithTrx(trxHandle *gorm.DB) TrxService
	List(c *gin.Context, userID uint) ([]models.Trx, error)
	Create(trxRequest *models.TrxRequest, userID uint) error
	Patch(transaction models.TrxPatchRequest, userID uint) error
	Delete(c *gin.Context, userID uint) error
}
