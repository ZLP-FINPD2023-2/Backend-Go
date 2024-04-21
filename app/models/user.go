package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Gender string

const (
	Female Gender = "Female"
	Male   Gender = "Male"
)

// Структура ответа на GET запрос
type UserGetResponse struct {
	Email      *string `json:"email"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Patronymic string  `json:"patronymic"`
	Gender     Gender  `json:"gender"`
	Birthday   string  `json:"birthday"`
}

type User struct {
	gorm.Model
	Email      *string `gorm:"unique"`
	Password   string
	FirstName  string
	LastName   string
	Patronymic string
	Gender     Gender
	Birthday   time.Time
}

// TableName gives table name of model
func (u User) TableName() string {
	return "users"
}

func (user *User) BeforeSave(tx *gorm.DB) error {
	// Хэширование пароля
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)

	return nil
}
