package repository

import (
	"finapp/lib"
	"finapp/models"

	"gorm.io/gorm"
)

type GeneratorRepository struct {
	logger   lib.Logger
	database lib.Database
}

func NewGeneratorRepository(
	logger lib.Logger,
	database lib.Database,
) GeneratorRepository {
	return GeneratorRepository{
		logger:   logger,
		database: database,
	}
}

func (r GeneratorRepository) WithTrx(trxHandle *gorm.DB) GeneratorRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.database.DB = trxHandle
	return r
}

func (r GeneratorRepository) Store(generator *models.Generator) error {
	return r.database.Create(&generator).Error
}

func (r GeneratorRepository) List(userID uint) ([]models.Generator, error) {
	var generators []models.Generator
	err := r.database.Where("user_id = ?", userID).Find(&generators).Error
	if err != nil {
		return nil, err
	}

	return generators, nil
}

func (r GeneratorRepository) Get(id, userID uint) (models.Generator, error) {
	var generator models.Generator
	err := r.database.Where("id = ? AND user_id = ?", id, userID).First(&generator).Error
	if err != nil {
		return models.Generator{}, err
	}

	return generator, nil
}

func (r GeneratorRepository) Update(generator models.Generator, id, userID uint) (models.Generator, error) {
	var genResponse models.Generator
	if err := r.database.Model(&genResponse).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(&generator).Error; err != nil {
		return models.Generator{}, err
	}

	if err := r.database.Where("id = ? AND user_id = ?", id, userID).
		First(&genResponse).Error; err != nil {
		return models.Generator{}, err
	}

	return genResponse, nil
}

func (r GeneratorRepository) Delete(id, userID uint) error {
	return r.database.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Generator{}).Error
}
