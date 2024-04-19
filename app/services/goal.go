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

type GoalService struct {
	logger           lib.Logger
	repository       repository.GoalRepository
	budgetRepository repository.BudgetRepository
}

func NewGoalService(
	logger lib.Logger,
	repository repository.GoalRepository,
) domains.GoalService {
	return GoalService{
		logger:     logger,
		repository: repository,
	}
}

func (s GoalService) WithTrx(trxHandle *gorm.DB) domains.GoalService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

func (s GoalService) List(c *gin.Context, userID uint) ([]models.GoalCalcResponse, error) {
	var date time.Time
	if dateToStr := c.Query("date"); dateToStr != "" {
		dateToTemp, err := time.Parse(constants.DateFormat, dateToStr)
		if err != nil {
			return nil, err
		}
		date = dateToTemp
	}

	goals, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	var resp []models.GoalCalcResponse
	for _, goal := range goals {
		budgets, err := s.budgetRepository.ListOfGoal(userID, goal.ID)
		if err != nil {
			return nil, err
		}

		calc := models.GoalCalcResponse{
			ID:           goal.ID,
			Title:        goal.Title,
			TargetAmount: goal.TargetAmount,
			Amount:       make(map[time.Time]decimal.Decimal),
		}

		for _, v := range budgets {
			amount, err := s.budgetRepository.GetBudgetAmount(v.ID, userID, date)
			if err != nil {
				return nil, err
			}

			calc.Amount[date] = calc.Amount[date].Add(amount)
		}
		resp = append(resp, calc)
	}

	return resp, err
}

func (s GoalService) Get(c *gin.Context, userID uint) (models.GoalCalcResponse, error) {
	queryID := c.Param("id")
	if queryID == "" {
		return models.GoalCalcResponse{}, errors.New("goal id does not exists")
	}
	id, err := strconv.Atoi(queryID)
	if err != nil {
		return models.GoalCalcResponse{}, err
	}

	var date time.Time
	if dateToStr := c.Query("date"); dateToStr != "" {
		dateToTemp, err := time.Parse(constants.DateFormat, dateToStr)
		if err != nil {
			return models.GoalCalcResponse{}, err
		}
		date = dateToTemp
	}

	goal, err := s.repository.Get(uint(id), userID)
	if err != nil {
		return models.GoalCalcResponse{}, err
	}

	budgets, err := s.budgetRepository.ListOfGoal(userID, goal.ID)
	if err != nil {
		return models.GoalCalcResponse{}, err
	}

	resp := models.GoalCalcResponse{
		ID:           goal.ID,
		Title:        goal.Title,
		TargetAmount: goal.TargetAmount,
		Amount:       make(map[time.Time]decimal.Decimal),
	}

	for _, v := range budgets {
		amount, err := s.budgetRepository.GetBudgetAmount(v.ID, userID, date)
		if err != nil {
			return models.GoalCalcResponse{}, err
		}

		resp.Amount[date] = resp.Amount[date].Add(amount)
	}

	return models.GoalCalcResponse{}, nil
}

func (s GoalService) Store(request *models.GoalStoreRequest, userID uint) (models.GoalResponse, error) {
	targetAmount, err := decimal.NewFromString(request.TargetAmount)
	if err != nil {
		return models.GoalResponse{}, err
	}
	goal := models.Goal{
		UserID:       userID,
		TargetAmount: targetAmount,
		Title:        request.Title,
	}

	err = s.repository.Create(&goal)
	if err != nil {
		return models.GoalResponse{}, err
	}

	resp := models.GoalResponse{
		ID:           goal.ID,
		Title:        goal.Title,
		TargetAmount: goal.TargetAmount,
	}

	return resp, nil
}

func (s GoalService) Update(c *gin.Context, req models.GoalUpdateRequest, userID uint) (models.GoalResponse, error) {
	queryID := c.Param("id")
	if queryID == "" {
		return models.GoalResponse{}, errors.New("goal id does not exists")
	}
	id, err := strconv.Atoi(queryID)
	if err != nil {
		return models.GoalResponse{}, err
	}

	var targetAmount decimal.Decimal
	if req.TargetAmount != "" {
		amount, err := decimal.NewFromString(req.TargetAmount)
		if err != nil {
			return models.GoalResponse{}, err
		}
		targetAmount = amount
	}

	goal := models.Goal{
		Title:        req.Title,
		TargetAmount: targetAmount,
	}

	updateGoal, err := s.repository.Patch(goal, uint(id), userID)
	if err != nil {
		return models.GoalResponse{}, nil
	}

	resp := models.GoalResponse{
		ID:           updateGoal.ID,
		Title:        updateGoal.Title,
		TargetAmount: updateGoal.TargetAmount,
	}
	return resp, nil
}

func (s GoalService) Delete(c *gin.Context, UserID uint) error {
	queryID := c.Param("id")
	if queryID == "" {
		return errors.New("goal id does not exists")
	}
	id, err := strconv.Atoi(queryID)
	if err != nil {
		return err
	}

	return s.repository.Delete(uint(id), UserID)
}
