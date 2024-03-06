package repository

import (
	"gorm.io/gorm"

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

func (r BudgetRepository) List(userID uint) ([]models.BudgetCalc, error) {
	var budgets []models.BudgetCalc
	rows, err := r.Database.Raw("SELECT "+
		"budgets.id,"+
		"budgets.created_at,"+
		"budgets.updated_at,"+
		"budgets.deleted_at,"+
		"budgets.user_id,"+
		"budgets.title,"+
		"budgets.goal,"+
		"COALESCE(SUM(t1.amount), 0) - COALESCE(SUM(t2.amount), 0) AS amount "+
		"FROM "+
		"budgets "+
		"LEFT JOIN "+
		"(SELECT budget_to, SUM(amount) as amount FROM transactions GROUP BY budget_to) t1 "+
		"ON budgets.id = t1.budget_to "+
		"LEFT JOIN "+
		"(SELECT budget_from, SUM(amount) as amount FROM transactions GROUP BY budget_from) t2 "+
		"ON budgets.id = t2.budget_from "+
		"WHERE "+
		"budgets.user_id = ? "+
		"GROUP BY "+
		"budgets.id", userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp models.BudgetCalc
		if err := r.Database.ScanRows(rows, &temp); err != nil {
			return nil, err
		}
		budgets = append(budgets, temp)
	}

	return budgets, err
}

func (r BudgetRepository) Get(id, userID uint) (models.Budget, error) {
	var budget models.Budget
	err := r.Database.Where("user_id = ?", userID).Where("id = ?", id).First(&budget).Error
	return budget, err
}
