package domains

import (
	"finapp/models"

	"github.com/gin-gonic/gin"
)

type GeneratorService interface {
	Store(generator models.GeneratorStoreRequest, userID uint) (models.GeneratorResponse, error)
	List(userID uint) ([]models.GeneratorResponse, error)
	Get(c *gin.Context, userID uint) (models.GeneratorResponse, error)
	Update(c *gin.Context, generator models.GeneratorPatchRequest, userID uint) (models.GeneratorResponse, error)
	Delete(c *gin.Context, userID uint) error
}
