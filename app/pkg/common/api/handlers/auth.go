package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"app/pkg/common/config"
	"app/pkg/common/db"
	"app/pkg/common/models"
)

// Генерация токена
func GenerateToken(ID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	return token.SignedString([]byte(config.Cfg.SecretKey))
}

// Вход

// Структура запроса
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @summary Login
// @tags auth
// @Description Вход пользователя
// @ID login
// @Accept json
// @Produce json
// @Param req body loginRequest true "Данные пользователя"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	// TODO: Сделать для ошибок/успеха свои структуры

	// Парсинг запроса
	var req loginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request body"},
		)
		return
	}

	// Поиск пользователя
	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid email or password"},
		)
		return
	}

	// Сравнение хэша и пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid email or password"},
		)
		return
	}

	// Генерация JWT токена
	token, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)
		return
	}

	// Отправка токена
	c.JSON(
		http.StatusOK,
		gin.H{"token": token},
	)
}

// Регистрация

// Структура запроса
type registerRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        uint8  `json:"age"`
	Gender     bool   `json:"gender"`
}

// @summary Register
// @tags auth
// @Description Регистрация пользователя
// @ID register
// @Accept json
// @Produce json
// @Param user body registerRequest true "Данные пользователя"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	// TODO: Сделать для ошибок/успеха свои структуры

	// Парсинг запроса
	var registerRequest registerRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Invalid request body",
			},
		)
		return
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to hash password"},
		)
		return
	}
	registerRequest.Password = string(hashedPassword)

	// Сохранение пользователя в БД
	user := models.User{
		Email:      registerRequest.Email,
		Password:   registerRequest.Password,
		FirstName:  registerRequest.FirstName,
		LastName:   registerRequest.LastName,
		Patronymic: registerRequest.Patronymic,
		Age:        registerRequest.Age,
		Gender:     registerRequest.Gender,
	}
	// TODO: улучшить обработку ошибок
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Failed to safe user",
			},
		)
		return
	}

	// Отправка ответа
	c.JSON(
		http.StatusOK,
		gin.H{"message": "User registered successfully"},
	)
}

// Выход
func Logout(c *gin.Context) {
	// TODO: Реализовать выход
	// TODO: Сделать для ошибок/успеха свои структуры
	c.JSON(
		http.StatusInternalServerError,
		gin.H{"error": "Not realized"},
	)
}