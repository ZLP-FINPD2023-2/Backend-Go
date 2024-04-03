package services

import (
	"errors"
	"strconv"
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
	}*/
	var (
		dateFrom time.Time
		dateTo   time.Time
	)

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFromTemp, err := time.Parse(constants.DateFormat, dateFromStr)
		if err != nil {
			return nil, err
		}
		dateFrom = dateFromTemp
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateToTemp, err := time.Parse(constants.DateFormat, dateToStr)
		if err != nil {
			return nil, err
		}
		dateTo = dateToTemp
	}

	// Выполнение запроса
	trxs, err := s.repository.List(userID, dateFrom, dateTo)

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

func (s TrxService) Patch(transaction models.TrxPatchRequest, userID uint) error {
	date, err := time.Parse(constants.DateFormat, transaction.Date)
	if err != nil {
		return err
	}

	amount, err := decimal.NewFromString(transaction.Amount)
	if err != nil {
		return err
	}

	trx := models.Trx{
		UserID:     userID,
		Title:      transaction.Title,
		Date:       date,
		Amount:     amount,
		BudgetFrom: transaction.BudgetFrom,
		BudgetTo:   transaction.BudgetTo,
	}
	trx.ID = transaction.ID

	return s.repository.Patch(trx)
}

func (s TrxService) Delete(c *gin.Context, userID uint) error {
	queryID := c.Query("id")
	if queryID == "" {
		return errors.New("trx id does not exists")
	}
	id, err := strconv.Atoi(queryID)
	if err != nil {
		return err
	}

	return s.repository.Delete(uint(id), userID)
}
