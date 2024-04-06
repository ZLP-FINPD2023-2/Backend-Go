package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TrxRequest struct {
	Title string `json:"title"`
	// Из-за datetime все может пойти по одному месту
	Date string `json:"date" validate:"required,datetime,isNotFutureDate"`
	// Из-за numeric все может пойти по одному месту
	Amount     string `json:"amount" validate:"required,numeric"`
	BudgetFrom uint   `json:"from" validate:"required"`
	BudgetTo   uint   `json:"to" validate:"required"`
}

type TrxResponse struct {
	ID         uint            `json:"id"`
	Title      string          `json:"title"`
	Date       time.Time       `json:"date"`
	Amount     decimal.Decimal `json:"amount"`
	BudgetFrom uint            `json:"from"`
	BudgetTo   uint            `json:"to"`
}

type TrxPatchRequest struct {
	ID         uint   `json:"id" validate:"required"`
	Title      string `json:"title"`
	Date       string `json:"date" validate:"datetime"`
	Amount     string `json:"amount" validate:"numeric"`
	BudgetFrom uint   `json:"from"`
	BudgetTo   uint   `json:"to"`
}

type Trx struct {
	gorm.Model
	UserID     uint
	Title      string
	Date       time.Time
	Amount     decimal.Decimal `sql:"type:decimal(20,2);"`
	BudgetFrom uint
	BudgetTo   uint
}

func (t Trx) TableName() string {
	return "transactions"
}
