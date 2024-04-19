package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"

	"finapp/lib"
	"finapp/models"
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
	var amountStr string
	err := r.Database.Model(&models.Trx{}).
		Select("SUM(CASE WHEN budget_to = ? THEN amount ELSE 0 END) - "+
			"SUM(CASE WHEN budget_from = ? THEN amount ELSE 0 END)", budgetID, budgetID).
		Where("user_id = ? AND date <= ?", userID, date).
		Group("user_id").
		Row().
		Scan(&amountStr)
	if err != nil {
		return decimal.Decimal{}, err
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return amount, nil
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
