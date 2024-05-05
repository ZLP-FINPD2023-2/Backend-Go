package repository

import (
	"errors"
	"time"

	"finapp/lib"
	"finapp/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BudgetRepository struct {
	logger   lib.Logger
	Database lib.Database
}

func NewBudgetRepository(logger lib.Logger, db lib.Database) BudgetRepository {
	return BudgetRepository{
		logger:   logger,
		Database: db,
	}
}

func (r BudgetRepository) WithTrx(trxHandle *gorm.DB) BudgetRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

func (r BudgetRepository) List(userID uint) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.Database.Where("user_id = ?", userID).Find(&budgets).Error
	if err != nil {
		return nil, err
	}
	return budgets, err
}

func (r BudgetRepository) ListOfGoal(userID uint, goalID uint) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.Database.Where("user_id = ? AND goal_id = ?", userID, goalID).Find(&budgets).Error
	if err != nil {
		return nil, err
	}
	return budgets, err
}

func (r BudgetRepository) Get(id uint, userID uint) (models.Budget, error) {
	var budget models.Budget
	err := r.Database.Where("user_id = ?", userID).Where("id = ?", id).First(&budget).Error
	return budget, err
}

// Получает сумму бюджета до определенной даты
func (r BudgetRepository) GetBudgetAmount(budgetID, userID uint, date time.Time) (decimal.Decimal, error) {
	var amount decimal.Decimal
	err := r.Database.Model(&models.Trx{}).
		Select("SUM(CASE WHEN budget_to = ? THEN CAST(amount AS DECIMAL) ELSE 0 END) - "+
			"SUM(CASE WHEN budget_from = ? THEN CAST(amount AS DECIMAL) ELSE 0 END)", budgetID, budgetID).
		Where("user_id = ? AND date <= ?", userID, date).
		Group("user_id").
		Row().
		Scan(&amount)
	if err != nil {
		return decimal.Decimal{}, err
	}

	var (
		genTo   models.Generator
		genFrom models.Generator
	)
	if err := r.Database.Where("user_id = ? AND budget_to = ?",
		userID,
		budgetID).First(&genTo).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return decimal.Decimal{}, err
	}

	if !genTo.DateFrom.IsZero() {
		var (
			currDate = genTo.DateFrom
			lastDate time.Time
			dayAdd   int
			monthAdd int
			yearAdd  int
		)

		switch genTo.Periodicity {
		case models.PeriodicityDaily:
			dayAdd, monthAdd, yearAdd = int(genTo.PeriodicityFactor), 0, 0
		case models.PeriodicityMonthly:
			dayAdd, monthAdd, yearAdd = 0, int(genTo.PeriodicityFactor), 0
		case models.PeriodicityYearly:
			dayAdd, monthAdd, yearAdd = 0, 0, int(genTo.PeriodicityFactor)
		}

		for currDate.Before(lastDate) || currDate.Equal(lastDate) {
			amount = amount.Add(genTo.Amount)
			currDate = currDate.AddDate(yearAdd, monthAdd, dayAdd)
		}

	}

	if err := r.Database.Where("user_id = ? AND budget_from = ?",
		userID,
		budgetID).First(&genFrom).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return decimal.Decimal{}, err
	}
	if !genFrom.DateFrom.IsZero() {
		var (
			currDate = genFrom.DateFrom
			lastDate time.Time
			dayAdd   int
			monthAdd int
			yearAdd  int
		)

		switch genFrom.Periodicity {
		case models.PeriodicityDaily:
			dayAdd, monthAdd, yearAdd = int(genFrom.PeriodicityFactor), 0, 0
		case models.PeriodicityMonthly:
			dayAdd, monthAdd, yearAdd = 0, int(genFrom.PeriodicityFactor), 0
		case models.PeriodicityYearly:
			dayAdd, monthAdd, yearAdd = 0, 0, int(genFrom.PeriodicityFactor)
		}

		if genFrom.DateTo.Time.IsZero() || genFrom.DateTo.Time.After(date) {
			lastDate = date
		} else {
			lastDate = genFrom.DateTo.Time
		}

		for currDate.Before(lastDate) || currDate.Equal(lastDate) {
			amount = amount.Sub(genFrom.Amount)
			currDate = currDate.AddDate(yearAdd, monthAdd, dayAdd)
		}
	}

	return amount, nil
}

func (r TrxRepository) GetBudgetChanges(budgetID, userID uint, dateFrom, dateTo time.Time) ([]models.BudgetChanges, error) {
	var changes []models.BudgetChanges
	query := r.Database.Model(&models.Trx{}).Select("SUM(CASE WHEN budget_to = ? THEN CAST(amount AS DECIMAL) ELSE 0 END) - "+
		"SUM(CASE WHEN budget_from = ? THEN CAST(amount AS DECIMAL) ELSE 0 END) as amount_change, date", budgetID, budgetID).
		Where("user_id = ?", userID).
		Where("budget_to = ? or budget_from = ?", budgetID, budgetID).
		Where("date > ?", dateFrom)
	if !dateTo.IsZero() {
		query.Where("date <= ?", dateTo)
	}
	if err := query.Group("date").Find(&changes).Error; err != nil {
		return nil, err
	}

	if dateTo.IsZero() {
		for _, v := range changes {
			if v.Date.After(dateTo) {
				dateTo = v.Date
			}
		}
	}

	var (
		genTo   []models.Generator
		genFrom []models.Generator
	)
	if err := r.Database.Where("user_id = ? AND budget_to = ?",
		userID,
		budgetID).Find(&genTo).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err := r.Database.Where("user_id = ? AND budget_from = ?",
		userID,
		budgetID).Find(&genFrom).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	for _, gen := range genTo {
		var (
			currDate = gen.DateFrom
			lastDate time.Time
			dayAdd   int
			monthAdd int
			yearAdd  int
		)

		if !gen.DateTo.Time.IsZero() && dateTo.After(gen.DateTo.Time) {
			lastDate = gen.DateTo.Time
		} else {
			lastDate = dateTo
		}

		switch gen.Periodicity {
		case models.PeriodicityDaily:
			dayAdd = int(gen.PeriodicityFactor)
		case models.PeriodicityMonthly:
			monthAdd = int(gen.PeriodicityFactor)
		case models.PeriodicityYearly:
			yearAdd = int(gen.PeriodicityFactor)
		}

		for currDate.Before(dateFrom) {
			currDate = currDate.AddDate(yearAdd, monthAdd, dayAdd)
		}

		for currDate.Before(lastDate) || currDate.Equal(lastDate) {
			changes = append(changes, models.BudgetChanges{AmountChange: gen.Amount, Date: currDate})
			currDate = currDate.AddDate(yearAdd, monthAdd, dayAdd)
		}
	}

	for _, gen := range genFrom {
		var (
			currDate = gen.DateFrom
			lastDate time.Time
			dayAdd   int
			monthAdd int
			yearAdd  int
		)

		if !gen.DateTo.Time.IsZero() && dateTo.After(gen.DateTo.Time) {
			lastDate = gen.DateTo.Time
		} else {
			lastDate = dateTo
		}

		switch gen.Periodicity {
		case models.PeriodicityDaily:
			dayAdd = int(gen.PeriodicityFactor)
		case models.PeriodicityMonthly:
			monthAdd = int(gen.PeriodicityFactor)
		case models.PeriodicityYearly:
			yearAdd = int(gen.PeriodicityFactor)
		}

		for currDate.Before(dateFrom) {
			currDate = currDate.AddDate(yearAdd, monthAdd, dayAdd)
		}

		for currDate.Before(lastDate) || currDate.Equal(lastDate) {
			changes = append(changes, models.BudgetChanges{AmountChange: gen.Amount.Neg(), Date: currDate})
			currDate = currDate.AddDate(yearAdd, monthAdd, dayAdd)
		}
	}

	return changes, nil
}

func (r BudgetRepository) Create(budget *models.Budget) error {
	return r.Database.Create(&budget).Error
}

func (r BudgetRepository) Patch(budget *models.Budget, id, userID uint) (models.Budget, error) {
	var budgetResponse models.Budget
	err := r.Database.Model(&budgetResponse).Where("user_id = ? AND id = ?", userID, id).Updates(&budget).Error
	if err != nil {
		return models.Budget{}, nil
	}

	if err := r.Database.Where("id = ? AND user_id = ?", id, userID).First(&budgetResponse).Error; err != nil {
		return models.Budget{}, err
	}
	return budgetResponse, nil
}

func (r BudgetRepository) Delete(id uint, userID uint) error {
	return r.Database.Where("user_id = ?", userID).Delete(&models.Budget{}, id).Error
}
