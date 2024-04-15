package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"log"
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

func (r TrxRepository) List(userID uint, dateFrom, dateTo time.Time, minAmount, maxAmount decimal.Decimal) ([]models.Trx, error) {
	var trxs []models.Trx
	query := r.Database.Where("user_id = ?", userID)
	if !dateFrom.Equal(time.Time{}) {
		query = query.Where("date >= ?", dateFrom)
		log.Println("Z A H O D 1")
	}
	if !dateTo.Equal(time.Time{}) {
		query = query.Where("date <= ?", dateTo)
		log.Println("Z A H O D 2")
	}
	if !maxAmount.Equal(decimal.Zero) {
		query = query.Where("amount <= ?", maxAmount)
		log.Println("Z A H O D 3")
	}
	query = query.Where("amount >= ?", minAmount)
	err := query.Find(&trxs).Error
	return trxs, err
}

func (r TrxRepository) ListFromBudget(budgetID, userID uint, dateFrom time.Time, dateTo time.Time) ([]models.Trx, error) {
	var trxs []models.Trx
	query := r.Database.Where("user_id = ?", userID)
	query = query.Where("date > ? AND date <= ?", dateFrom, dateTo)
	query = query.Where("budget_from = ?", budgetID).Or("budget_to = ?", budgetID)
	err := query.Find(&trxs).Error
	return trxs, err
}

func (r TrxRepository) Patch(trx models.Trx, id, userID uint) (models.Trx, error) {
	var trxResponse models.Trx
	if err := r.Database.Model(&trxResponse).Where("id = ? AND user_id = ?", id, userID).
		Updates(&trx).Error; err != nil {
		return models.Trx{}, err
	}

	if err := r.Database.Where("id = ? AND user_id = ?", id, userID).First(&trxResponse).Error; err != nil {
		return models.Trx{}, err
	}

	return trxResponse, nil
}

func (r TrxRepository) Delete(id uint, userID uint) error {
	return r.Database.Where("user_id = ?", userID).Delete(&models.Trx{}, id).Error
}
