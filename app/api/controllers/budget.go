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

type BudgetController struct {
	logger  lib.Logger
	service domains.BudgetService
}

func NewBudgetController(
	logger lib.Logger,
	service domains.BudgetService,
) BudgetController {
	return BudgetController{
		logger:  logger,
		service: service,
	}
}

// Получение

// @Deprecated
// @Security		ApiKeyAuth
// @summary		Get budgets
// @tags			budget
// @Description	Получение бюджетов
// @ID				budget-get
// @Accept			json
// @Produce		json
// @Success		200	{array}	models.BudgetGetResponse
// @Router			/budget [get]
func (bc BudgetController) List(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	budgets, err := bc.service.List(userID.(uint))
	// TODO: Улучшить обработку ошибок
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to get budgets",
			"description": err.Error(),
		})
		return
	}

	var budgetResponses []models.BudgetGetResponse
	for _, budget := range budgets {
		budgetResponses = append(budgetResponses, models.BudgetGetResponse{
			Title:  budget.Title,
			ID:     budget.ID,
			Amount: budget.Amount,
		})
	}

	c.JSON(http.StatusOK, budgetResponses)
}

// Создание

// @Deprecated
// @Security		ApiKeyAuth
// @summary		Create budget
// @tags			budget
// @Description	Создание бюджета
// @ID				budget-create
// @Accept			json
// @Produce		json
// @Param			budget	body	models.BudgetCreateRequest	true	"Данные бюждета"
// @Router			/budget [post]
func (bc BudgetController) Post(c *gin.Context) {
	var budget models.BudgetCreateRequest

	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(budget); err != nil {
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

	err := bc.service.Create(&budget, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Budget added successfully",
	})
}

// Обновление

// @Deprecated
// @Security		ApiKeyAuth
// @summary		Patch budget
// @tags			budget
// @Description	Изменение бюджета
// @ID				budget-patch
// @Accept			json
// @Produce		json
// @Param			budget	body	models.BudgetPatchRequest	true	"Данные бюждета"
// @Router			/budget [patch]
func (bc BudgetController) Patch(c *gin.Context) {
	var budget models.BudgetPatchRequest

	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(budget); err != nil {
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

	if err := bc.service.Patch(budget, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to update budget: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "budget was updated",
	})
}

// Удаление

// @Deprecated
// @Security		ApiKeyAuth
// @summary		Delete budget
// @tags			budget
// @Description	Удаление бюджета
// @ID				budget-delete
// @Accept			json
// @Produce		json
// @Param			id	query	integer	false	"id бюджета"
// @Router			/budget [delete]
func (bc BudgetController) Delete(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	if err := bc.service.Delete(c, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to delete budget: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "budget deleted successfully",
	})
}
