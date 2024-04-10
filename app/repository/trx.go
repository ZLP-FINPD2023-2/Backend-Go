package repository

import (
	"gorm.io/gorm"
	"time"

	"finapp/lib"
	"finapp/models"
)

type TrxRepository struct {
	logger   lib.Logger
	Database lib.Database
}

func NewTrxRepository(
	logger lib.Logger,
	db lib.Database,
) TrxRepository {
	return TrxRepository{
		logger:   logger,
		Database: db,
	}
}

func (r TrxRepository) WithTrx(trxHandle *gorm.DB) TrxRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

func (r TrxRepository) Create(model *models.Trx) error {
	return r.Database.Create(&model).Error
}

func (r TrxRepository) List(userID uint, dateFrom time.Time, dateTo time.Time) ([]models.Trx, error) {
	var trxs []models.Trx
	query := r.Database.Where("user_id = ?", userID)
	query = query.Where("date >= ?", dateFrom)
	query = query.Where("date <= ?", dateTo)
	err := query.Find(&trxs).Error
	return trxs, err
}

func (r TrxRepository) ListFromBudget(budgetID, userID uint, dateFrom time.Time, dateTo time.Time) ([]models.Trx, error) {
	var trxs []models.Trx
	query := r.Database.Where("user_id = ?", userID)
	query = query.Where("date > ?", dateFrom)
	query = query.Where("date <= ?", dateTo)
	query.Where("budget_from = ? OR budget_to = ?", budgetID, budgetID)
	err := query.Find(&trxs).Error
	return trxs, err
}

func (r TrxRepository) Patch(trx models.Trx) error {
	return r.Database.Save(&trx).Error
}

func (r TrxRepository) Delete(id uint, userID uint) error {
	return r.Database.Where("user_id = ?", userID).Delete(&models.Trx{}, id).Error
}
