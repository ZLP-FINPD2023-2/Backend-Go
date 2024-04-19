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

// @Deprecated
// @Security ApiKeyAuth
// @summary List goals
// @tags goal
// @Description Получение бюджетов
// @ID goal-list
// @Accept json
// @Produce json
// @Param date query string false "Дата в формате 18-10-2004"
// @Success 200 {array} models.GoalResponse
// @Router /goal [get]
func (gc GoalController) List(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	goals, err := gc.service.List(c, userID.(uint))
	// TODO: Улучшить обработку ошибок
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to get goals",
			"description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// @Deprecated
// @Security ApiKeyAuth
// @summary List goals
// @tags goal
// @Description Получение бюджетов
// @ID goal-get
// @Accept json
// @Produce json
// @Param id path integer false "id цели"
// @Param date query string false "Дата в формате 18-10-2004"
// @Success 200 {object} models.GoalResponse
// @Router /goal [get]
func (gc GoalController) Get(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	goal, err := gc.service.Get(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get goal: %s", err.Error()),
		})
	}

	c.JSON(http.StatusOK, goal)
}

// Создание

// @Security ApiKeyAuth
// @summary Store goal
// @tags goal
// @Description Создание цели
// @ID goal-create
// @Accept json
// @Produce json
// @Param goal body models.GoalStoreRequest true "Данные бюждета"
// @Success 200 {object} models.GoalResponse
// @Router /goal [post]
func (gc GoalController) Store(c *gin.Context) {
	var goal models.GoalStoreRequest

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

	resp, err := gc.service.Store(&goal, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Обновление

// @Security ApiKeyAuth
// @summary Update goal
// @tags goal
// @Description Изменение цели
// @ID goal-patch
// @Accept json
// @Produce json
// @Param id path integer false "id цели"
// @Param goal body models.GoalUpdateRequest true "Данные цели"
// @Success 200 {object} models.GoalResponse
// @Router /goal/{id} [patch]
func (gc GoalController) Update(c *gin.Context) {
	var goal models.GoalUpdateRequest

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

	resp, err := gc.service.Update(c, goal, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to update goal: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Удаление

// @Security ApiKeyAuth
// @summary Delete goal
// @tags goal
// @Description Удаление цели
// @ID goal-delete
// @Accept json
// @Produce json
// @Param id path integer false "id цели"
// @Router /goal/{id} [delete]
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

	c.JSON(http.StatusNoContent, gin.H{
		"message": "goal was deleted",
	})
}
