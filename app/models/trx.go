package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TrxRequest struct {
	Title      string `json:"title"`
	Date       string `json:"date" validate:"required,isNotFutureDate"`
	Amount     string `json:"amount" validate:"required,numeric"`
	BudgetFrom *uint  `json:"from"`
	BudgetTo   *uint  `json:"to"`
}

type TrxResponse struct {
	ID         uint            `json:"id"`
	Title      string          `json:"title"`
	Date       string          `json:"date"`
	Amount     decimal.Decimal `json:"amount"`
	BudgetFrom *uint           `json:"from"`
	BudgetTo   *uint           `json:"to"`
}

type TrxPatchRequest struct {
	Title  string `json:"title"`
	Amount string `json:"amount"`
}

type Trx struct {
	gorm.Model
	UserID     uint
	Title      string
	Date       time.Time
	Amount     decimal.Decimal `sql:"type:decimal(20,2);"`
	BudgetFrom *sql.NullInt64
	BudgetTo   *sql.NullInt64
}

func (t Trx) TableName() string {
	return "transactions"
}
