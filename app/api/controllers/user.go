package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"finapp/constants"
	"finapp/domains"
	"finapp/lib"
	"finapp/models"
)

// UserController struct
type UserController struct {
	logger  lib.Logger
	service domains.UserService
}

// NewUserController creates new controller
func NewUserController(
	logger lib.Logger,
	service domains.UserService,
) UserController {
	return UserController{
		logger:  logger,
		service: service,
	}
}

// Удаление

// @Security ApiKeyAuth
// @summary Delete user
// @tags user
// @Description Удаление пользователя
// @ID delete_user
// @Accept json
// @Produce json
// @Router /user [delete]
func (uc UserController) Delete(c *gin.Context) {
	// Парсинг запроса
	userId, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	// Удаление пользователя
	if err := uc.service.Delete(userId.(uint)); err != nil {
		// Необработанные ошибки
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	// Отправка ответа
	c.JSON(http.StatusNoContent, gin.H{
		"message": "User deleted successfully",
	})
}

// Получение

// @Security ApiKeyAuth
// @summary Get user
// @tags user
// @Description Получение пользователя
// @ID get_user
// @Accept json
// @Produce json
// @Router /user [get]
func (uc UserController) Get(c *gin.Context) {
	userID, ok := c.Get(constants.UserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	user, err := uc.service.Get(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	response := models.UserGetResponse{
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Patronymic: user.Patronymic,
		Gender:     user.Gender,
		Birthday:   user.Birthday.Format(constants.DateFormat),
	}

	c.JSON(http.StatusOK, response)
}

// Обновление

// @Deprecated
// @Security ApiKeyAuth
// @summary Update user
// @tags user
// @Description Обновление пользователя
// @ID update_user
// @Accept json
// @Produce json
// @Router /user [patch]
func (uc UserController) Update(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
