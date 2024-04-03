package models

import (
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

type GoalCreateRequest struct {
	Title        string `json:"title" validate:"required"`
	TargetAmount string `json:"target_amount" validate:"required,numeric"`
}

type GoalGetResponse struct {
	ID           uint            `json:"id"`
	Title        string          `json:"title"`
	Amount       decimal.Decimal `json:"amount"`
	TargetAmount decimal.Decimal `json:"target_amount"`
}

type GoalPatchRequest struct {
	ID           uint   `json:"id" validate:"required"`
	Title        string `json:"title"`
	TargetAmount string `json:"target_amount"`
}

type Goal struct {
	gorm.Model
	UserID       uint
	Title        string          `gorm:"unique"`
	TargetAmount decimal.Decimal `sql:"type:decimal(20,2);"`
}

type GoalCalc struct {
	gorm.Model
	UserID       uint
	Title        string          `gorm:"unique"`
	Amount       decimal.Decimal `sql:"type:decimal(20,2);"`
	TargetAmount decimal.Decimal `sql:"type:decimal(20,2);"`
}

func (m Goal) TableName() string {
	return "goals"
}
