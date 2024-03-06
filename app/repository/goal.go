package repository

import (
	"gorm.io/gorm"

	"finapp/lib"
	"finapp/models"
)

type GoalRepository struct {
	logger   lib.Logger
	Database lib.Database
}

func NewGoalRepository(logger lib.Logger, db lib.Database) GoalRepository {
	return GoalRepository{
		logger:   logger,
		Database: db,
	}
}

func (r GoalRepository) WithTrx(trxHandle *gorm.DB) GoalRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

func (r GoalRepository) Get(id, userID uint) (models.Goal, error) {
	var goal models.Goal
	err := r.Database.Where("user_id = ?", userID).Where("id = ?", id).First(&goal).Error
	return goal, err
}

func (r GoalRepository) List(userID uint) ([]models.GoalCalc, error) {
	var goals []models.GoalCalc
	rows, err := r.Database.Raw("SELECT "+
		"goals.id,"+
		"goals.created_at,"+
		"goals.updated_at,"+
		"goals.deleted_at,"+
		"goals.user_id,"+
		"goals.title,"+
		"COALESCE(SUM(t3.amount), 0) AS total_amount "+
		"FROM goals "+
		"JOIN "+
		"(SELECT "+
		"budgets.id,"+
		"budgets.created_at,"+
		"budgets.updated_at,"+
		"budgets.deleted_at,"+
		"budgets.user_id,"+
		"budgets.title,"+
		"budgets.goal,"+
		"COALESCE(SUM(t1.amount), 0) - COALESCE(SUM(t2.amount), 0) AS amount "+
		"FROM budgets "+
		"LEFT JOIN "+
		"(SELECT "+
		"budget_to, "+
		"SUM(amount) as amount "+
		"FROM transactions GROUP BY budget_to) t1 ON budgets.id = t1.budget_to "+
		"LEFT JOIN "+
		"(SELECT "+
		"budget_from, "+
		"SUM(amount) as amount "+
		"FROM transactions GROUP BY budget_from) t2 ON budgets.id = t2.budget_from "+
		"WHERE "+
		"budgets.user_id = ? GROUP BY budgets.id) t3 ON goals.id = t3.goal GROUP BY goals.id", userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp models.GoalCalc
		if err := r.Database.ScanRows(rows, &temp); err != nil {
			return nil, err
		}
		goals = append(goals, temp)
	}
	return goals, err
}
