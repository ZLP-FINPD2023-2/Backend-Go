package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"finapp/domains"
	"finapp/lib"
)

// UserController struct
type UserController struct {
	logger      lib.Logger
	service     domains.AuthService
	userService domains.UserService
}

// NewUserController creates new controller
func NewUserController(
	logger lib.Logger,
	service domains.AuthService,
	userService domains.UserService,
) UserController {
	return UserController{
		logger:      logger,
		service:     service,
		userService: userService,
	}
}

// Удаление

// @summary		Delete user
// @tags			user
// @Description	Удаление пользователя
// @ID				delete
// @Accept			json
// @Produce		json
// @Router			/user [delete]
func (uc UserController) Delete(c *gin.Context) {
	// Парсинг запроса
	userId, exists := c.Get("userId")
	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Удаление пользователя
	if err := uc.userService.Delete(uint(userId.(int))); err != nil {
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
