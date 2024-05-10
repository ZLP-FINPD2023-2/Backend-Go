package services

import (
	"database/sql"
	"errors"
	"sort"
	"strconv"
	"time"

	"finapp/constants"
	"finapp/domains"
	"finapp/lib"
	"finapp/models"
	"finapp/repository"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GoalService struct {
	logger           lib.Logger
	repository       repository.GoalRepository
	budgetRepository repository.BudgetRepository
	trxRepository    repository.TrxRepository
}

func NewGoalService(
	logger lib.Logger,
	repository repository.GoalRepository,
	budgetRepository repository.BudgetRepository,
	trxRepository repository.TrxRepository,
) domains.GoalService {
	return GoalService{
		logger:           logger,
		repository:       repository,
		budgetRepository: budgetRepository,
		trxRepository:    trxRepository,
	}
}

func (s GoalService) WithTrx(trxHandle *gorm.DB) domains.GoalService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

func (s GoalService) List(c *gin.Context, userID uint) ([]models.GoalCalcResponse, error) {
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

		g := models.GoalCalcResponse{
			ID:           goal.ID,
			Title:        goal.Title,
			TargetAmount: goal.TargetAmount.InexactFloat64(),
			Amounts:      make(map[string]float64),
		}

		changes := make(map[time.Time]decimal.Decimal)
		for _, v := range budgets {
			amount, err := s.budgetRepository.GetBudgetAmount(v.ID, userID, dateFrom)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return nil, err
				}
				amount = decimal.Zero
			}
			if !dateFrom.IsZero() {
				g.Amounts[dateFrom.Format(constants.DateFormat)] =
					g.Amounts[dateFrom.Format(constants.DateFormat)] + amount.InexactFloat64()
			}

			budgetChanges, err := s.trxRepository.GetBudgetChanges(v.ID, userID, dateFrom, dateTo)
			if err != nil {
				return nil, err
			}

			for _, change := range budgetChanges {
				changes[change.Date] = changes[change.Date].Add(change.AmountChange)
			}
		}

		dates := make([]time.Time, 0, len(changes))
		for k, _ := range changes {
			dates = append(dates, k)
		}
		sort.Slice(dates, func(i int, j int) bool {
			return dates[i].Before(dates[j])
		})

		var (
			currDate        time.Time
			currAmountState = g.Amounts[dateFrom.Format(constants.DateFormat)]
		)
		if !dateFrom.IsZero() {
			currDate = dateFrom
		}

		for _, v := range dates {
			if !currDate.IsZero() {
				for currDate.Before(v) {
					currDate = currDate.Add(24 * time.Hour)
					g.Amounts[currDate.Format(constants.DateFormat)] =
						g.Amounts[currDate.Format(constants.DateFormat)] + currAmountState
				}
			}
			currAmountState = currAmountState + changes[v].InexactFloat64()
			g.Amounts[currDate.Format(constants.DateFormat)] = currAmountState
			currDate = v
		}

		if !dateTo.IsZero() {
			for !currDate.Equal(dateTo) && currDate.Before(dateTo) {
				currDate = currDate.Add(24 * time.Hour)
				g.Amounts[currDate.Format(constants.DateFormat)] = currAmountState
			}
		}

		resp = append(resp, g)
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

	var (
		dateFrom time.Time
		dateTo   time.Time
	)
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFromTemp, err := time.Parse(constants.DateFormat, dateFromStr)
		if err != nil {
			return models.GoalCalcResponse{}, err
		}
		dateFrom = dateFromTemp
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateToTemp, err := time.Parse(constants.DateFormat, dateToStr)
		if err != nil {
			return models.GoalCalcResponse{}, err
		}
		dateTo = dateToTemp
	}

	if !dateTo.IsZero() && dateTo.Before(dateFrom) {
		return models.GoalCalcResponse{}, errors.New("date_from time goes after date_to")
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
		TargetAmount: goal.TargetAmount.InexactFloat64(),
		Amounts:      make(map[string]float64),
	}

	changes := make(map[time.Time]decimal.Decimal)
	for _, v := range budgets {
		amount, err := s.budgetRepository.GetBudgetAmount(v.ID, userID, dateFrom)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return models.GoalCalcResponse{}, err
			}
			amount = decimal.Zero
		}
		if !dateFrom.IsZero() {
			resp.Amounts[dateFrom.Format(constants.DateFormat)] =
				resp.Amounts[dateFrom.Format(constants.DateFormat)] + amount.InexactFloat64()
		}

		budgetChanges, err := s.trxRepository.GetBudgetChanges(v.ID, userID, dateFrom, dateTo)
		if err != nil {
			return models.GoalCalcResponse{}, err
		}

		for _, change := range budgetChanges {
			changes[change.Date] = changes[change.Date].Add(change.AmountChange)
		}
	}

	dates := make([]time.Time, 0, len(changes))
	for k, _ := range changes {
		dates = append(dates, k)
	}
	sort.Slice(dates, func(i int, j int) bool {
		return dates[i].Before(dates[j])
	})

	var (
		currDate        time.Time
		currAmountState = resp.Amounts[dateFrom.Format(constants.DateFormat)]
	)
	if !dateFrom.IsZero() {
		currDate = dateFrom
	}

	for _, v := range dates {
		if !currDate.IsZero() {
			for currDate.Before(v) {
				currDate = currDate.Add(24 * time.Hour)
				resp.Amounts[currDate.Format(constants.DateFormat)] =
					resp.Amounts[currDate.Format(constants.DateFormat)] + currAmountState
			}
		}
		currAmountState = currAmountState + changes[v].InexactFloat64()
		resp.Amounts[currDate.Format(constants.DateFormat)] = currAmountState
		currDate = v
	}

	if !dateTo.IsZero() {
		for !currDate.Equal(dateTo) && currDate.Before(dateTo) {
			currDate = currDate.Add(24 * time.Hour)
			resp.Amounts[currDate.Format(constants.DateFormat)] = currAmountState
		}
	}

	return resp, nil
}

func (s GoalService) Store(request *models.GoalStoreRequest, userID uint) (models.GoalResponse, error) {
	targetAmount := decimal.NewFromFloat(request.TargetAmount)
	goal := models.Goal{
		UserID:       userID,
		TargetAmount: targetAmount,
		Title:        request.Title,
	}

	err := s.repository.Create(&goal)
	if err != nil {
		return models.GoalResponse{}, err
	}

	resp := models.GoalResponse{
		ID:           goal.ID,
		Title:        goal.Title,
		TargetAmount: goal.TargetAmount.InexactFloat64(),
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
	if req.TargetAmount != 0 {
		amount := decimal.NewFromFloat(req.TargetAmount)
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
		TargetAmount: updateGoal.TargetAmount.InexactFloat64(),
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
