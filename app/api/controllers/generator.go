package controllers

import (
	"fmt"
	"net/http"

	"finapp/constants"
	"finapp/domains"
	"finapp/lib"
	"finapp/lib/validators"
	"finapp/models"

	"github.com/gin-gonic/gin"
)

type GeneratorController struct {
	logger  lib.Logger
	service domains.GeneratorService
}

func NewGeneratorController(
	logger lib.Logger,
	service domains.GeneratorService,
) GeneratorController {
	return GeneratorController{
		logger:  logger,
		service: service,
	}
}

// @Security ApiKeyAuth
// @summary Create generator
// @tags generator
// @Description Создание генератора транзакций
// @ID post_gen
// @Accept json
// @Produce json
// @Param transaction body models.GeneratorStoreRequest true "Данные генератора"
// @Success 200 {object} models.GeneratorResponse
// @Router /trx/generator [post]
func (gc GeneratorController) Store(c *gin.Context) {
	var generator models.GeneratorStoreRequest

	if err := c.ShouldBindJSON(&generator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(generator); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": validators.ParseValidationErrors(err),
		})
		return
	}

	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	resp, err := gc.service.Store(generator, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to store generator: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @summary List of generator
// @tags generator
// @Description Получение генераторов транзакций
// @ID list_gen
// @Accept json
// @Produce json
// @Success 200 {array} models.GeneratorResponse
// @Router /trx/generator [get]
func (gc GeneratorController) List(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	resp, err := gc.service.List(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get list of generators: %s", err.Error()),
		})
		return
	}

	if resp == nil {
		resp = make([]models.GeneratorResponse, 0)
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @summary Get generator
// @tags generator
// @Description Получение генератора транзакций
// @ID get_gen
// @Accept json
// @Produce json
// @Param  id  path  int  true  "ID генератора"
// @Success 200 {object} models.GeneratorResponse
// @Router /trx/generator/{id} [get]
func (gc GeneratorController) Get(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	resp, err := gc.service.Get(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get generator: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @summary Update generator
// @tags generator
// @Description Изменение генератора транзакций
// @ID patch_gen
// @Accept json
// @Produce json
// @Param  id  path  int  true  "ID генератора"
// @Param transaction body models.GeneratorPatchRequest true "Данные генератора"
// @Success 200 {object} models.GeneratorResponse
// @Router /trx/generator/{id} [patch]
func (gc GeneratorController) Update(c *gin.Context) {
	var generator models.GeneratorPatchRequest

	if err := c.ShouldBindJSON(&generator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(generator); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": validators.ParseValidationErrors(err),
		})
		return
	}

	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	resp, err := gc.service.Update(c, generator, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to update generator: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @summary Delete generator
// @tags generator
// @Description Удаление генератора транзакций
// @ID delete_gen
// @Accept json
// @Produce json
// @Param  id  path  int  true  "ID генератора"
// @Router /trx/generator/{id} [delete]
func (gc GeneratorController) Delete(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	if err := gc.service.Delete(c, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to delete generator: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "transaction was deleted",
	})
}
