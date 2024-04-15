package services

import (
	"errors"
	"finapp/constants"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"strconv"
	"time"

	"finapp/domains"
	"finapp/lib"
	"finapp/models"
	"finapp/repository"
)

type BudgetService struct {
	logger         lib.Logger
	repository     repository.BudgetRepository
	goalRepository repository.GoalRepository
	trxRepository  repository.TrxRepository
}

func NewBudgetService(
	logger lib.Logger,
	repository repository.BudgetRepository,
	goalRepository repository.GoalRepository,
	trxRepository repository.TrxRepository,
) domains.BudgetService {
	return BudgetService{
		logger:         logger,
		repository:     repository,
		goalRepository: goalRepository,
		trxRepository:  trxRepository,
	}
}

func (s BudgetService) WithTrx(trxHandle *gorm.DB) domains.BudgetService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

func (s BudgetService) Get(c *gin.Context, userID uint) (models.BudgetGetResponse, error) {
	paramID := c.Params.ByName("id")
	if paramID == "" {
		return models.BudgetGetResponse{}, errors.New("budget id does not exists")
	}
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return models.BudgetGetResponse{}, err
	}

	var (
		dateFrom time.Time
		dateTo   time.Time
	)

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFromTemp, err := time.Parse(constants.DateFormat, dateFromStr)
		if err != nil {
			return models.BudgetGetResponse{}, err
		}
		dateFrom = dateFromTemp
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateToTemp, err := time.Parse(constants.DateFormat, dateToStr)
		if err != nil {
			return models.BudgetGetResponse{}, err
		}
		dateTo = dateToTemp
	}

	budget, err := s.repository.Get(uint(id), userID)
	if err != nil {
		return models.BudgetGetResponse{}, err
	}

	startAmount, err := s.repository.GetBudgetAmount(uint(id), userID, dateFrom)
	if err != nil {
		return models.BudgetGetResponse{}, err
	}

	trxs, err := s.trxRepository.ListFromBudget(uint(id), userID, dateFrom, dateTo)
	if err != nil {
		return models.BudgetGetResponse{}, err
	}

	budgetCalc := models.BudgetGetResponse{
		ID:      budget.ID,
		Title:   budget.Title,
		Goal:    budget.Goal,
		Amounts: make(map[time.Time]decimal.Decimal),
	}
	currAmount := startAmount
	budgetCalc.Amounts[dateFrom] = currAmount
	for _, v := range trxs {
		if v.BudgetFrom.Int64 == int64(id) {
			currAmount = currAmount.Sub(v.Amount)
		} else {
			currAmount = currAmount.Add(v.Amount)
		}
		budgetCalc.Amounts[v.Date] = currAmount
	}

	return budgetCalc, nil
}

func (s BudgetService) List(c *gin.Context, userID uint) ([]models.BudgetGetResponse, error) {
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

	budgets, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	var budgetsAmounts []models.BudgetGetResponse
	for _, v := range budgets {
		budg := models.BudgetGetResponse{
			ID:      v.ID,
			Title:   v.Title,
			Goal:    v.Goal,
			Amounts: make(map[time.Time]decimal.Decimal),
		}

		startAmount, err := s.repository.GetBudgetAmount(v.ID, userID, dateFrom)
		if err != nil {
			return nil, err
		}

		currAmount := startAmount
		budg.Amounts[dateFrom] = startAmount

		trxs, err := s.trxRepository.ListFromBudget(v.ID, userID, dateFrom, dateTo)
		if err != nil {
			return nil, err
		}

		for _, trx := range trxs {
			if trx.BudgetFrom.Int64 == int64(v.ID) {
				currAmount = currAmount.Sub(trx.Amount)
			} else {
				currAmount = currAmount.Add(trx.Amount)
			}
			budg.Amounts[trx.Date] = currAmount
		}

		budgetsAmounts = append(budgetsAmounts, budg)
	}
	return budgetsAmounts, err
}

func (s BudgetService) Create(request *models.BudgetCreateRequest, userID uint) (models.BudgetCreateResponse, error) {
	if request.Goal != 0 {
		_, err := s.goalRepository.Get(request.Goal, userID)
		if err != nil {
			return models.BudgetCreateResponse{}, err
		}
	}

	budget := models.Budget{
		UserID: userID,
		Title:  request.Title,
		Goal:   request.Goal,
	}

	if err := s.repository.Create(&budget); err != nil {
		return models.BudgetCreateResponse{}, err
	}

	newBudget := models.BudgetCreateResponse{
		ID:     budget.ID,
		Title:  budget.Title,
		GoadID: budget.Goal,
	}

	return newBudget, nil
}

func (s BudgetService) Patch(budget models.BudgetPatchRequest, userID uint) (models.BudgetPatchResponse, error) {
	updateBudget := models.Budget{
		UserID: userID,
		Title:  budget.Title,
		Goal:   budget.Goal,
	}
	updateBudget.ID = budget.ID

	err := s.repository.Patch(&updateBudget)
	if err != nil {
		return models.BudgetPatchResponse{}, err
	}

	resp := models.BudgetPatchResponse{
		ID:    updateBudget.ID,
		Title: updateBudget.Title,
		Goal:  updateBudget.Goal,
	}

	return resp, nil
}

func (s BudgetService) Delete(c *gin.Context, userID uint) error {
	paramID := c.Params.ByName("id")
	if paramID == "" {
		return errors.New("budget id does not exists")
	}
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return err
	}

	return s.repository.Delete(uint(id), userID)
}
