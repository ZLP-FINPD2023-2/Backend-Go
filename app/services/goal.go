package services

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"strconv"

	"finapp/domains"
	"finapp/lib"
	"finapp/models"
	"finapp/repository"
)

type GoalService struct {
	logger     lib.Logger
	repository repository.GoalRepository
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

func (s GoalService) List(userID uint) ([]models.GoalCalc, error) {
	goals, err := s.repository.List(userID)

	return goals, err
}

func (s GoalService) Create(request *models.GoalCreateRequest, userID uint) error {
	targetAmount, err := decimal.NewFromString(request.TargetAmount)
	if err != nil {
		return err
	}
	goal := models.Goal{
		UserID:       userID,
		TargetAmount: targetAmount,
		Title:        request.Title,
	}

	return s.repository.Create(goal)
}

func (s GoalService) Patch(req models.GoalPatchRequest, userID uint) error {
	targetAmount, err := decimal.NewFromString(req.TargetAmount)
	if err != nil {
		return err
	}
	goal := models.Goal{
		UserID:       userID,
		Title:        req.Title,
		TargetAmount: targetAmount,
	}
	goal.ID = req.ID

	return s.repository.Patch(goal)
}

func (s GoalService) Delete(c *gin.Context, UserID uint) error {
	queryID := c.Query("id")
	if queryID == "" {
		return errors.New("goal id does not exists")
	}
	id, err := strconv.Atoi(queryID)
	if err != nil {
		return err
	}
	return s.repository.Delete(uint(id), UserID)
}
