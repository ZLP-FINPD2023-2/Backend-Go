package models

import (
	"finapp/lib/validators"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

type Goal struct {
	gorm.Model
	UserID uint   `validate:"required"`
	Title  string `validate:"required" gorm:"unique"`
}

type GoalCalc struct {
	gorm.Model
	UserID uint            `validate:"required"`
	Title  string          `validate:"required" gorm:"unique"`
	Amount decimal.Decimal `validate:"required" sql:"type:decimal(20,2);"`
}

func (m Goal) TableName() string {
	return "goals"
}

func (m *Goal) BeforeSave(db *gorm.DB) error {
	// Валидация
	validate := validators.CustomValidator()
	if err := validate.Struct(m); err != nil {
		return err
	}

	return nil
}

type GoalCreateRequest struct {
	Title string `json:"title"`
}

type GoalGetResponse struct {
	Title  string          `json:"title"`
	ID     uint            `json:"id"`
	Amount decimal.Decimal `json:"amount"`
}
