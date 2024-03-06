package services

import (
	"time"

	"gorm.io/gorm"

	"finapp/constants"
	"finapp/domains"
	"finapp/lib"
	"finapp/models"
	"finapp/repository"
)

// UserService service layer
type UserService struct {
	logger     lib.Logger
	repository repository.UserRepository
}

// NewUserService creates a new userservice
func NewUserService(logger lib.Logger, repository repository.UserRepository) domains.UserService {
	return UserService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s UserService) WithTrx(trxHandle *gorm.DB) domains.UserService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// Register call to register the user
func (s UserService) Register(q *models.RegisterRequest) error {
	var err error
	birthday, err := time.Parse(constants.DateFormat, q.Birthday)
	if err != nil {
		return err
	}

	user := models.User{
		Email:      q.Email,
		Password:   q.Password,
		FirstName:  q.FirstName,
		LastName:   q.LastName,
		Patronymic: q.Patronymic,
		Gender:     q.Gender,
		Birthday:   birthday,
	}

	return s.repository.Create(&user)
}

func (s UserService) GetUserByEmail(email *string) (*models.User, error) {
	return s.repository.GetByEmail(email)
}

func (s UserService) Get(id uint) (*models.User, error) {
	return s.repository.Get(id)
}

// UpdateUser updates the user
//func (s UserService) UpdateUser(user models.User) error {
//	return s.repository.Save(&user).Error
//}

// Delete deletes the user
func (s UserService) Delete(id uint) error {
	return s.repository.Delete(id)
}
