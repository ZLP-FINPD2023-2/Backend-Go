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

type GoalController struct {
	logger  lib.Logger
	service domains.GoalService
}

func NewGoalController(
	logger lib.Logger,
	service domains.GoalService,
) GoalController {
	return GoalController{
		logger:  logger,
		service: service,
	}
}

// Получение

// @Security		ApiKeyAuth
// @summary		List goals
// @tags			goal
// @Description	Получение бюджетов
// @ID				goal-list
// @Accept			json
// @Produce		json
// @Success		200	{array}	models.GoalGetResponse
// @Router			/goal [get]
func (gc GoalController) List(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	goals, err := gc.service.List(userID.(uint))
	// TODO: Улучшить обработку ошибок
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to get goals",
			"description": err.Error(),
		})
		return
	}

	var goalResponses []models.GoalGetResponse
	for _, goal := range goals {
		goalResponses = append(goalResponses, models.GoalGetResponse{
			Title:  goal.Title,
			ID:     goal.ID,
			Amount: goal.Amount,
		})
	}

	c.JSON(http.StatusOK, goalResponses)
}

// Создание

// @Security		ApiKeyAuth
// @summary		Create goal
// @tags			goal
// @Description	Создание цели
// @ID				goal-create
// @Accept			json
// @Produce		json
// @Param			goal	body	models.GoalCreateRequest	true	"Данные бюждета"
// @Router			/goal [post]
func (gc GoalController) Create(c *gin.Context) {
	var goal models.GoalCreateRequest

	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(goal); err != nil {
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

	err := gc.service.Create(&goal, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Goal added successfully",
	})
}

// @Security		ApiKeyAuth
// @summary		Patch goal
// @tags			goal
// @Description	Изменение цели
// @ID				goal-patch
// @Accept			json
// @Produce		json
// @Param			goal	body	models.GoalPatchRequest	true	"Данные цели"
// @Router			/goal [patch]
func (gc GoalController) Patch(c *gin.Context) {
	var goal models.GoalPatchRequest

	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(goal); err != nil {
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

	if err := gc.service.Patch(goal, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to update goal: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "goal was updated",
	})
}

// @Security		ApiKeyAuth
// @summary		Delete goal
// @tags			goal
// @Description	Удаление цели
// @ID				goal-delete
// @Accept			json
// @Produce		json
// @Param			id	query	integer	false	"id бюджета"
// @Router			/goal [delete]
func (gc GoalController) Delete(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	if err := gc.service.Delete(c, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to delete goal: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "goal was deleted",
	})
}
