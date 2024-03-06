package services

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"finapp/constants"
	"finapp/domains"
	"finapp/lib"
	"finapp/models"
	"finapp/repository"
)

type TrxService struct {
	logger           lib.Logger
	repository       repository.TrxRepository
	budgetRepository repository.BudgetRepository
}

func NewTrxService(
	logger lib.Logger,
	repository repository.TrxRepository,
	budgetRepository repository.BudgetRepository,
) domains.TrxService {
	return TrxService{
		logger:           logger,
		repository:       repository,
		budgetRepository: budgetRepository,
	}
}

func (s TrxService) WithTrx(trxHandle *gorm.DB) domains.TrxService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

func (s TrxService) List(c *gin.Context, userID uint) ([]models.Trx, error) {
	var trxs []models.Trx

	// Создание запроса
	query := s.repository.Database.Where("user_id = ?", userID)

	// Фильтрация запроса
	/* TODO: Реализовать фильтрацию по сумме

	if amountMinStr := c.Query("amount_min"); amountMinStr != "" {
		amountMin, err := decimal.NewFromString(amountMinStr)
		if err != nil {
			return nil, err
		}
		query = query.Where("amount >= ?", amountMin)
	}

	if amountMaxStr := c.Query("amount_max"); amountMaxStr != "" {
		amountMax, err := decimal.NewFromString(amountMaxStr)
		if err != nil {
			return nil, err
		}
		query = query.Where("amount <= ?", amountMax)
	}
	*/

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFrom, err := time.Parse(constants.DateFormat, dateFromStr)
		if err != nil {
			return nil, err
		}
		query = query.Where("date >= ?", dateFrom)
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateTo, err := time.Parse(constants.DateFormat, dateToStr)
		if err != nil {
			return nil, err
		}
		query = query.Where("date <= ?", dateTo)
	}

	// Выполнение запроса
	err := query.Find(&trxs).Error

	return trxs, err
}

func (s TrxService) Create(trxRequest *models.TrxRequest, userID uint) error {
	date, err := time.Parse(constants.DateFormat, trxRequest.Date)
	if err != nil {
		return err
	}

	amount, err := decimal.NewFromString(trxRequest.Amount)
	if err != nil {
		return err
	}

	_, err = s.budgetRepository.Get(trxRequest.BudgetTo, userID)
	if err != nil {
		return err
	}

	_, err = s.budgetRepository.Get(trxRequest.BudgetFrom, userID)
	if err != nil {
		return err
	}

	transaction := models.Trx{
		UserID:     userID,
		Title:      trxRequest.Title,
		Date:       date,
		Amount:     amount,
		BudgetTo:   trxRequest.BudgetTo,
		BudgetFrom: trxRequest.BudgetFrom,
	}

	return s.repository.Create(&transaction)
}
