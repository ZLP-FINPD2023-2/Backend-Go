package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TrxRequest struct {
	Title      string  `json:"title"`
	Date       string  `json:"date" validate:"required,isNotFutureDate"`
	Amount     float64 `json:"amount" validate:"required,numeric"`
	BudgetFrom *uint   `json:"budget_from"`
	BudgetTo   *uint   `json:"budget_to"`
}

type TrxResponse struct {
	ID         uint    `json:"id"`
	Title      string  `json:"title"`
	Date       string  `json:"date"`
	Amount     float64 `json:"amount"`
	BudgetFrom *uint   `json:"budget_from"`
	BudgetTo   *uint   `json:"budget_to"`
}

type TrxPatchRequest struct {
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
}

type Trx struct {
	gorm.Model
	UserID          uint
	User            User `gorm:"foreignKey:UserID"`
	Title           string
	Date            time.Time
	Amount          decimal.Decimal `sql:"type:decimal(20,2);"`
	BudgetFrom      *sql.NullInt64
	BudgetToModel   Budget `gorm:"foreignKey:BudgetTo"`
	BudgetTo        *sql.NullInt64
	BudgetFromModel Budget `gorm:"foreignKey:BudgetFrom"`
}

func (t Trx) TableName() string {
	return "transactions"
}
