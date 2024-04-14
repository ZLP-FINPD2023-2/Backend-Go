package controllers

import (
	"finapp/lib/validators"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"finapp/constants"
	"finapp/domains"
	"finapp/lib"
	"finapp/models"
)

type TrxController struct {
	logger  lib.Logger
	service domains.TrxService
}

func NewTrxController(
	logger lib.Logger,
	service domains.TrxService,
) TrxController {
	return TrxController{
		logger:  logger,
		service: service,
	}
}

// Получение

// @Deprecated
// @Security ApiKeyAuth
// @summary List trx
// @tags trx
// @Description Получение транзакции
// @ID get_trx
// @Accept json
// @Produce json
// @Param date_from query string false "Дата начала периода в формате 18-10-2004"
// @Param date_to query string false "Дата окончания периода в формате 18-10-2004"
// @Success 200 {array} models.TrxResponse
// @Router /trx [get]
func (tc TrxController) List(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	trxs, err := tc.service.List(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to get transaction",
			"description": err.Error(),
		})
		return
	}

	var trxResponses []models.TrxResponse
	for _, trx := range trxs {
		trxResponses = append(trxResponses, models.TrxResponse{
			ID:         trx.ID,
			Title:      trx.Title,
			Date:       trx.Date,
			Amount:     trx.Amount,
			BudgetFrom: trx.BudgetFrom,
			BudgetTo:   trx.BudgetTo,
		})
	}

	c.JSON(http.StatusOK, trxResponses)
}

// Создание

// @Deprecated
// @Security ApiKeyAuth
// @summary Create trx
// @tags trx
// @Description Создание транзакции
// @ID post
// @Accept json
// @Produce json
// @Param transaction body models.TrxRequest true "Данные пользователя"
// @Router /trx [post]
func (tc TrxController) Post(c *gin.Context) {
	var transaction models.TrxRequest

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(transaction); err != nil {
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

	err := tc.service.Create(&transaction, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction added successfully",
	})
}

// Обновление

// @Deprecated
// @Security ApiKeyAuth
// @summary Patch trx
// @tags trx
// @Description Изменение транзакции
// @ID post
// @Accept json
// @Produce json
// @Param transaction body models.TrxPatchRequest true "Данные транзакций"
// @Router /trx [patch]
func (tc TrxController) Patch(c *gin.Context) {
	var transaction models.TrxPatchRequest
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(transaction); err != nil {
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

	if err := tc.service.Patch(transaction, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to update trx: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transaction was updated",
	})
}

// Удаление

// @Deprecated
// @Security ApiKeyAuth
// @summary Delete trx
// @tags trx
// @Description Удаление транзакции
// @ID post
// @Accept json
// @Produce json
// @Param id query integer false "id транзакции"
// @Router /trx [delete]
func (tc TrxController) Delete(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	if err := tc.service.Delete(c, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to delete trx: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transaction was deleted",
	})
}
