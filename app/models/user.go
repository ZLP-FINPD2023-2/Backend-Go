package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"finapp/lib/validators"
)

type Gender string

const (
	Female Gender = "Female"
	Male   Gender = "Male"
)

type User struct {
	gorm.Model
	Email      *string `gorm:"unique" validate:"required,email"`
	Password   string  `validate:"required,min=8"`
	FirstName  string  `validate:"required"`
	LastName   string  `validate:"required"`
	Patronymic string
	Gender     Gender    `validate:"oneof=Male Female"`
	Birthday   time.Time `validate:"required"`
}

// TableName gives table name of model
func (u User) TableName() string {
	return "users"
}

func (user *User) BeforeSave(tx *gorm.DB) error {
	// Валидация
	validate := validators.CustomValidator()
	if err := validate.Struct(user); err != nil {
		return err
	}

	// Хэширование пароля
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashPassword)

	return nil
}

// Структура ответа на GET запрос
type GetResponse struct {
	Email      *string `json:"email"`
	First_name string  `json:"first_name"`
	Last_name  string  `json:"last_name"`
	Patronymic string  `json:"patronymic"`
	Gender     Gender  `json:"gender"`
	Birthday   string  `json:"birthday"`
}
