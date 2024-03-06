package services

import (
	"gorm.io/gorm"

	"finapp/domains"
	"finapp/lib"
	"finapp/models"
	"finapp/repository"
)

type BudgetService struct {
	logger         lib.Logger
	repository     repository.BudgetRepository
	goalRepository repository.GoalRepository
}

func NewBudgetService(
	logger lib.Logger,
	repository repository.BudgetRepository,
	goalRepository repository.GoalRepository,
) domains.BudgetService {
	return BudgetService{
		logger:         logger,
		repository:     repository,
		goalRepository: goalRepository,
	}
}

func (s BudgetService) WithTrx(trxHandle *gorm.DB) domains.BudgetService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

func (s BudgetService) List(userID uint) ([]models.BudgetCalc, error) {
	budgets, err := s.repository.List(userID)
	return budgets, err
}

func (s BudgetService) Create(request *models.BudgetCreateRequest, userID uint) error {
	if request.Goal != 0 {
		_, err := s.goalRepository.Get(request.Goal, userID)
		if err != nil {
			return err
		}
	}

	transaction := models.Budget{
		UserID: userID,
		Title:  request.Title,
		Goal:   request.Goal,
	}

	return s.repository.Database.Create(&transaction).Error
}
