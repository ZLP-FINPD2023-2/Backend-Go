package domains

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"finapp/models"
)

type TrxService interface {
	WithTrx(trxHandle *gorm.DB) TrxService
	List(c *gin.Context, userID uint) ([]models.TrxResponse, error)
	Create(trxRequest *models.TrxRequest, userID uint) (models.TrxResponse, error)
	Patch(c *gin.Context, transaction models.TrxPatchRequest, userID uint) (models.TrxResponse, error)
	Delete(c *gin.Context, userID uint) error
}
