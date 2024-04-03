package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BudgetCreateRequest struct {
	Title string `json:"title" validate:"required"`
	Goal  uint   `json:"goal" validate:"required"`
}

type BudgetPatchRequest struct {
	ID    uint   `json:"id" validate:"required"`
	Title string `json:"title"`
	Goal  uint   `json:"goal"`
}

type BudgetGetResponse struct {
	Title  string          `json:"title"`
	ID     uint            `json:"id"`
	Goal   uint            `json:"goal"`
	Amount decimal.Decimal `json:"amount"`
}

type Budget struct {
	gorm.Model
	UserID uint
	Title  string `gorm:"unique"`
	Goal   uint
}

type BudgetCalc struct {
	gorm.Model
	UserID uint
	Title  string `gorm:"unique"`
	Goal   uint
	Amount decimal.Decimal `sql:"type:decimal(20,2);"`
}

func (b Budget) TableName() string {
	return "budgets"
}
