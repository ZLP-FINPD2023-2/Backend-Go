package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"finapp/domains"
	"finapp/lib"
	"finapp/lib/validators"
	"finapp/models"
)

// JWTAuthController struct
type JWTAuthController struct {
	logger      lib.Logger
	service     domains.AuthService
	userService domains.UserService
}

// NewJWTAuthController creates new controller
func NewJWTAuthController(
	logger lib.Logger,
	service domains.AuthService,
	userService domains.UserService,
) JWTAuthController {
	return JWTAuthController{
		logger:      logger,
		service:     service,
		userService: userService,
	}
}

// Вход

// @summary		Login
// @tags			auth
// @Description	Вход пользователя
// @ID				login
// @Accept			json
// @Produce		json
// @Param			req	body	models.LoginRequest	true	"Данные пользователя"
// @Router			/auth/login [post]
func (jwt JWTAuthController) Login(c *gin.Context) {
	// Парсинг запроса
	var q models.LoginRequest
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(q); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": validators.ParseValidationErrors(err),
		})
		return
	}

	// Нахождение пользователя по email пользователя
	user, err := jwt.userService.GetUserByEmail(q.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Сравнение хэша и пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(q.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Получение токена
	token, err := jwt.service.CreateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Отправка токена
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Регистрация

// @summary		Register
// @tags			auth
// @Description	Регистрация пользователя
// @ID				register
// @Accept			json
// @Produce		json
// @Param			user	body	models.RegisterRequest	true	"Данные пользователя"
// @Router			/auth/register [post]
func (jwt JWTAuthController) Register(c *gin.Context) {
	// Парсинг запроса
	var q models.RegisterRequest
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := validators.IsValid(q); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": validators.ParseValidationErrors(err),
		})
		return
	}

	// Регистрация пользователя
	user, err := jwt.userService.Register(&q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register user",
		})
		return
	}

	// Отправка ответа
	c.JSON(http.StatusOK, user)
}
