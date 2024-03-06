package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"finapp/lib/validators"
)

type Trx struct {
	gorm.Model
	UserID     uint            `validate:"required"`
	Title      string          `validate:"required"`
	Date       time.Time       `validate:"required,isNotFutureDate"`
	Amount     decimal.Decimal `validate:"required" sql:"type:decimal(20,2);"`
	BudgetFrom uint            `validate:"required"`
	BudgetTo   uint            `validate:"required"`
}

func (t Trx) TableName() string {
	return "transactions"
}

func (trx *Trx) BeforeSave(db *gorm.DB) error {
	// Валидация
	validate := validators.CustomValidator()
	if err := validate.Struct(trx); err != nil {
		return err
	}

	return nil
}

type TrxResponse struct {
	Title      string          `json:"title"`
	Date       time.Time       `json:"date"`
	Amount     decimal.Decimal `json:"amount"`
	BudgetFrom uint            `validate:"required"`
	BudgetTo   uint            `validate:"required"`
}

type TrxRequest struct {
	Title      string `json:"title"`
	Date       string `json:"date"`
	Amount     string `json:"amount"`
	BudgetFrom uint   `json:"from" validate:"required"`
	BudgetTo   uint   `json:"to" validate:"required"`
}
