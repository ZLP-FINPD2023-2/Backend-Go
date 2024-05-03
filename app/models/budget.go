package models

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type BudgetCreateRequest struct {
	Title string `json:"title" validate:"required"`
	Goal  *uint  `json:"goal_id" validate:"required"`
}

type BudgetCreateResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	GoadID *uint  `json:"goad_id"`
}

type BudgetPatchRequest struct {
	Title string `json:"title"`
	Goal  *uint  `json:"goal_id"`
}

type BudgetPatchResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Goal  *uint  `json:"goal_id"`
}

type BudgetGetResponse struct {
	Title   string             `json:"title"`
	ID      uint               `json:"id"`
	Goal    *uint              `json:"goal_id"`
	Amounts map[string]float64 `json:"amounts"`
}

type Budget struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	Title  string
	GoalID *sql.NullInt64
	Goal   Goal `gorm:"foreignKey:GoalID"`
}

type BudgetChanges struct {
	AmountChange decimal.Decimal
	Date         time.Time
}

func (b Budget) TableName() string {
	return "budgets"
}
