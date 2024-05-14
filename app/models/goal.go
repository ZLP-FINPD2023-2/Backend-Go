package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// DB
// / Models
type Goal struct {
	gorm.Model
	UserID       uint
	User         User `gorm:"foreignKey:UserID"`
	Title        string
	TargetAmount decimal.Decimal `sql:"type:decimal(20,2);"`
}

// / Table Name
func (m Goal) TableName() string {
	return "goals"
}

// Requests/Responses
// / Store
type GoalStoreRequest struct {
	Title        string  `json:"title" validate:"required"`
	TargetAmount float64 `json:"target_amount" validate:"required,numeric"`
}

type GoalCalcResponse struct {
	ID           uint               `json:"id"`
	Title        string             `json:"title"`
	Amounts      map[string]float64 `json:"amount"`
	TargetAmount float64            `json:"target_amount"`
}

// / Get
type GoalResponse struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	TargetAmount float64 `json:"target_amount"`
}

// Update
type GoalUpdateRequest struct {
	Title        string  `json:"title"`
	TargetAmount float64 `json:"target_amount"`
}
