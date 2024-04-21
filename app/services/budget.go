package services

import (
	"database/sql"
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
	logger        lib.Logger
	repository    repository.BudgetRepository
	trxRepository repository.TrxRepository
}

func NewBudgetService(
	logger lib.Logger,
	repository repository.BudgetRepository,
	trxRepository repository.TrxRepository,
) domains.BudgetService {
	return BudgetService{
		logger:        logger,
		repository:    repository,
		trxRepository: trxRepository,
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

	if !dateTo.IsZero() && dateTo.Before(dateFrom) {
		return models.BudgetGetResponse{}, errors.New("date_from time goes after date_to")
	}

	budget, err := s.repository.Get(uint(id), userID)
	if err != nil {
		return models.BudgetGetResponse{}, err
	}

	startAmount, err := s.repository.GetBudgetAmount(uint(id), userID, dateFrom)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.BudgetGetResponse{}, err
		}
		startAmount = decimal.New(0, 0)
	}

	changes, err := s.trxRepository.GetBudgetChanges(uint(id), userID, dateFrom, dateTo)
	if err != nil {
		return models.BudgetGetResponse{}, err
	}

	resp := models.BudgetGetResponse{
		ID:      budget.ID,
		Goal:    budget.GoalID,
		Title:   budget.Title,
		Amounts: make(map[string]decimal.Decimal),
	}

	var currAmount decimal.Decimal
	if !dateFrom.IsZero() || !startAmount.Equal(decimal.Zero) {
		resp.Amounts[dateFrom.Format(constants.DateFormat)] = startAmount
		currAmount = startAmount
	}

	for _, change := range changes {
		currAmount = currAmount.Add(change.AmountChange)
		resp.Amounts[change.Date.Format(constants.DateFormat)] = currAmount
	}
	return resp, nil
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

	if !dateTo.IsZero() && dateTo.Before(dateFrom) {
		return nil, errors.New("date_from time goes after date_to")
	}

	budgets, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	var budgetsAmounts []models.BudgetGetResponse
	for _, v := range budgets {
		budget, err := s.repository.Get(v.ID, userID)
		if err != nil {
			return nil, err
		}

		startAmount, err := s.repository.GetBudgetAmount(v.ID, userID, dateFrom)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			startAmount = decimal.New(0, 0)
		}

		changes, err := s.trxRepository.GetBudgetChanges(v.ID, userID, dateTo, dateFrom)
		if err != nil {
			return nil, err
		}

		budg := models.BudgetGetResponse{
			ID:      budget.ID,
			Goal:    budget.GoalID,
			Title:   budget.Title,
			Amounts: make(map[string]decimal.Decimal),
		}

		var currAmount decimal.Decimal
		if !dateFrom.IsZero() || !startAmount.Equal(decimal.Zero) {
			budg.Amounts[dateFrom.Format(constants.DateFormat)] = startAmount
			currAmount = startAmount
		}

		for _, change := range changes {
			currAmount = currAmount.Add(change.AmountChange)
			budg.Amounts[change.Date.Format(constants.DateFormat)] = currAmount
		}

		budgetsAmounts = append(budgetsAmounts, budg)
	}

	return budgetsAmounts, err
}

func (s BudgetService) Create(request *models.BudgetCreateRequest, userID uint) (models.BudgetCreateResponse, error) {
	budget := models.Budget{
		UserID: userID,
		Title:  request.Title,
		GoalID: request.Goal,
	}

	if err := s.repository.Create(&budget); err != nil {
		return models.BudgetCreateResponse{}, err
	}

	newBudget := models.BudgetCreateResponse{
		ID:     budget.ID,
		Title:  budget.Title,
		GoadID: budget.GoalID,
	}

	return newBudget, nil
}

func (s BudgetService) Patch(c *gin.Context, budget models.BudgetPatchRequest, userID uint) (models.BudgetPatchResponse, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return models.BudgetPatchResponse{}, errors.New("id does not exists")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.BudgetPatchResponse{}, err
	}

	updateBudget := models.Budget{
		Title:  budget.Title,
		GoalID: budget.Goal,
	}

	budgetDB, err := s.repository.Patch(&updateBudget, uint(id), userID)
	if err != nil {
		return models.BudgetPatchResponse{}, err
	}

	resp := models.BudgetPatchResponse{
		ID:    budgetDB.ID,
		Title: budgetDB.Title,
		Goal:  budgetDB.GoalID,
	}

	return resp, nil
}

func (s BudgetService) Delete(c *gin.Context, userID uint) error {
	paramID := c.Param("id")
	if paramID == "" {
		return errors.New("budget id does not exists")
	}
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return err
	}

	return s.repository.Delete(uint(id), userID)
}
