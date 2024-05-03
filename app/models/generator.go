package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Periodicity string

const (
	PeriodicityDaily   Periodicity = "Daily"
	PeriodicityMonthly Periodicity = "Monthly"
	PeriodicityYearly  Periodicity = "Yearly"
)

type GeneratorStoreRequest struct {
	Title             string      `json:"title"`
	Amount            float64     `json:"amount" validate:"numeric"`
	Periodicity       Periodicity `json:"periodicity"`
	PeriodicityFactor uint        `json:"periodicity_factor"`
	BudgetFrom        *uint       `json:"budget_from"`
	BudgetTo          *uint       `json:"budget_to"`
	DateFrom          string      `json:"date_from"`
	DateTo            *string     `json:"date_to"`
}

type GeneratorPatchRequest struct {
	Title             string      `json:"title"`
	Amount            float64     `json:"amount" validate:"numeric"`
	Periodicity       Periodicity `json:"periodicity"`
	PeriodicityFactor uint        `json:"periodicity_factor"`
	BudgetFrom        *uint       `json:"budget_from"`
	BudgetTo          *uint       `json:"budget_to"`
	DateTo            string      `json:"date_to"`
	DateFrom          string      `json:"date_from"`
}

type GeneratorResponse struct {
	ID                uint        `json:"id"`
	Title             string      `json:"title"`
	Amount            float64     `json:"amount"`
	Periodicity       Periodicity `json:"periodicity"`
	PeriodicityFactor uint        `json:"periodicity_factor"`
	BudgetFrom        *uint       `json:"budget_from"`
	BudgetTo          *uint       `json:"budget_to"`
	DateFrom          string      `json:"date_from"`
	DateTo            *string     `json:"date_to"`
}

type Generator struct {
	gorm.Model
	UserID            uint
	User              User `gorm:"foreignKey:UserID"`
	Title             string
	Amount            decimal.Decimal
	Periodicity       Periodicity
	PeriodicityFactor uint
	BudgetFrom        *sql.NullInt64
	BudgetFromModel   Budget `gorm:"foreignKey:BudgetFrom"`
	BudgetTo          *sql.NullInt64
	BudgetToModel     Budget `gorm:"foreignKey:BudgetTo"`
	DateFrom          time.Time
	DateTo            *sql.NullTime
}
