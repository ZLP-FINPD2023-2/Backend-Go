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
	Title        string
	TargetAmount decimal.Decimal `sql:"type:decimal(20,2);"`
}

type GoalCalc struct {
	gorm.Model
	UserID       uint
	Title        string
	Amount       decimal.Decimal `sql:"type:decimal(20,2);"`
	TargetAmount decimal.Decimal `sql:"type:decimal(20,2);"`
}

// / Table Name
func (m Goal) TableName() string {
	return "goals"
}

// Requests/Responses
// / Store
type GoalStoreRequest struct {
	Title        string `json:"title" validate:"required"`
	TargetAmount string `json:"target_amount" validate:"required,numeric"`
}

// / Get
type GoalListResponse struct {
	ID           uint            `json:"id"`
	Title        string          `json:"title"`
	Amount       decimal.Decimal `json:"amount"`
	TargetAmount decimal.Decimal `json:"target_amount"`
}

// Update
type GoalUpdateRequest struct {
	ID           uint   `json:"id" validate:"required"`
	Title        string `json:"title"`
	TargetAmount string `json:"target_amount"`
}
