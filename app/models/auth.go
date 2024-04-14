package models

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// Register
type RegisterRequest struct {
	Email      *string `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required,min=8"`
	FirstName  string  `json:"firstname" validate:"required"`
	LastName   string  `json:"lastname" validate:"required"`
	Patronymic string  `json:"patronymic,omitempty"`
	Gender     Gender  `json:"gender" validate:"required,oneof=Male Female"`
	// Есть вероятность, что из-за datetime все полетит в одно место
	Birthday string `json:"birthday" validate:"required"`
}

type RegisterResponse struct {
	Email      string `json:"email"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Patronymic string `json:"patronymic"`
	Gender     Gender `json:"gender"`
}

// Login
type LoginRequest struct {
	Email    *string `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
}

type LoginResponce struct {
	Token string `json:"token"`
}
