package services

import (
	"database/sql"
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

func (s TrxService) List(c *gin.Context, userID uint) ([]models.TrxResponse, error) {
	var (
		dateFrom  time.Time
		dateTo    time.Time
		minAmount decimal.Decimal
		maxAmount decimal.Decimal
	)

	if amountMinStr := c.Query("amount_min"); amountMinStr != "" {
		amountMin, err := decimal.NewFromString(amountMinStr)
		if err != nil {
			return nil, err
		}
		minAmount = amountMin
	}
	if amountMaxStr := c.Query("amount_max"); amountMaxStr != "" {
		amountMax, err := decimal.NewFromString(amountMaxStr)
		if err != nil {
			return nil, err
		}
		maxAmount = amountMax
	}

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
	trxs, err := s.repository.List(userID, dateFrom, dateTo, minAmount, maxAmount)
	if err != nil {
		return nil, err
	}

	var trxResponses []models.TrxResponse
	for _, trx := range trxs {
		trxResponses = append(trxResponses, models.TrxResponse{
			ID:         trx.ID,
			Title:      trx.Title,
			Date:       trx.Date,
			Amount:     trx.Amount,
			BudgetFrom: convertBudgetID(trx.BudgetFrom),
			BudgetTo:   convertBudgetID(trx.BudgetTo),
		})
	}

	return trxResponses, nil
}

func (s TrxService) Create(trxRequest *models.TrxRequest, userID uint) (models.TrxResponse, error) {
	date, err := time.Parse(constants.DateFormat, trxRequest.Date)
	if err != nil {
		return models.TrxResponse{}, err
	}

	amount, err := decimal.NewFromString(trxRequest.Amount)
	if err != nil {
		return models.TrxResponse{}, err
	}

	// TODO: убрать этот позор, добавить foreign keys
	/*_, err = s.budgetRepository.Get(trxRequest.BudgetTo, userID)
	if err != nil {
		return err
	}

	_, err = s.budgetRepository.Get(trxRequest.BudgetFrom, userID)
	if err != nil {
		return err
	}*/

	transaction := models.Trx{
		UserID: userID,
		Title:  trxRequest.Title,
		Date:   date,
		Amount: amount,
		BudgetTo: func() *sql.NullInt64 {
			if trxRequest.BudgetTo != nil {
				return &sql.NullInt64{
					Int64: int64(*trxRequest.BudgetTo),
					Valid: true,
				}
			}
			return &sql.NullInt64{}
		}(),
		BudgetFrom: func() *sql.NullInt64 {
			if trxRequest.BudgetFrom != nil {
				return &sql.NullInt64{
					Int64: int64(*trxRequest.BudgetFrom),
					Valid: true,
				}
			}
			return &sql.NullInt64{}
		}(),
	}

	err = s.repository.Create(&transaction)
	if err != nil {
		return models.TrxResponse{}, err
	}

	trxResponse := models.TrxResponse{
		ID:         transaction.ID,
		Title:      transaction.Title,
		Date:       transaction.Date,
		Amount:     transaction.Amount,
		BudgetFrom: convertBudgetID(transaction.BudgetFrom),
		BudgetTo:   convertBudgetID(transaction.BudgetTo),
	}
	return trxResponse, nil
}

func (s TrxService) Patch(c *gin.Context, transaction models.TrxPatchRequest, userID uint) (models.TrxResponse, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return models.TrxResponse{}, errors.New("id doesn't exists")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.TrxResponse{}, err
	}

	var amount decimal.Decimal
	if transaction.Amount != "" {
		currAmount, err := decimal.NewFromString(transaction.Amount)
		if err != nil {
			return models.TrxResponse{}, err
		}
		amount = currAmount
	}

	trx := models.Trx{
		Title:  transaction.Title,
		Amount: amount,
	}

	trxUpdate, err := s.repository.Patch(trx, uint(id), userID)
	if err != nil {
		return models.TrxResponse{}, err
	}

	trxResponse := models.TrxResponse{
		ID:         trxUpdate.ID,
		Title:      trxUpdate.Title,
		Date:       trxUpdate.Date,
		Amount:     trxUpdate.Amount,
		BudgetFrom: convertBudgetID(trxUpdate.BudgetFrom),
		BudgetTo:   convertBudgetID(trxUpdate.BudgetTo),
	}

	return trxResponse, nil
}

func (s TrxService) Delete(c *gin.Context, userID uint) error {
	queryID := c.Param("id")
	if queryID == "" {
		return errors.New("trx id does not exists")
	}
	id, err := strconv.Atoi(queryID)
	if err != nil {
		return err
	}

	return s.repository.Delete(uint(id), userID)
}

func convertBudgetID(budget *sql.NullInt64) *uint {
	if budget == nil {
		return nil
	}
	var id uint
	if budget.Valid {
		id = uint(budget.Int64)
		return &id
	}
	return nil
}
