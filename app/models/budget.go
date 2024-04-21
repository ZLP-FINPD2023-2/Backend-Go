package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type BudgetCreateRequest struct {
	Title string `json:"title" validate:"required"`
	Goal  uint   `json:"goal_id" validate:"required"`
}

type BudgetCreateResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	GoadID uint   `json:"goad_id"`
}

type BudgetPatchRequest struct {
	Title string `json:"title"`
	Goal  uint   `json:"goal_id"`
}

type BudgetPatchResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Goal  uint   `json:"goal_id"`
}

type BudgetGetResponse struct {
	Title   string                     `json:"title"`
	ID      uint                       `json:"id"`
	Goal    uint                       `json:"goal"`
	Amounts map[string]decimal.Decimal `json:"amounts"`
}

type Budget struct {
	gorm.Model
	UserID uint
	Title  string
	GoalID uint
}

type BudgetChanges struct {
	AmountChange decimal.Decimal
	Date         time.Time
}

func (b Budget) TableName() string {
	return "budgets"
}
