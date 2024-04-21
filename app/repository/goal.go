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

func (r GoalRepository) List(userID uint) ([]models.Goal, error) {
	var goals []models.Goal
	if err := r.Database.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func (r GoalRepository) Create(goal *models.Goal) error {
	return r.Database.Create(&goal).Error
}

func (r GoalRepository) Patch(goal models.Goal, id, userID uint) (models.Goal, error) {
	var updateGoal models.Goal
	err := r.Database.Model(&updateGoal).Where("user_id = ? AND id = ?", userID, id).Updates(&goal).Error
	if err != nil {
		return models.Goal{}, nil
	}

	if err := r.Database.Where("id = ? AND user_id = ?", id, userID).First(&updateGoal).Error; err != nil {
		return models.Goal{}, err
	}
	return updateGoal, nil
}

func (r GoalRepository) Delete(id uint, userID uint) error {
	return r.Database.Where("user_id = ?", userID).Delete(&models.Goal{}, id).Error
}
