package services

import (
	"database/sql"
	"finapp/constants"
	"finapp/domains"
	"finapp/lib"
	"finapp/models"
	"finapp/repository"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type GeneratorService struct {
	logger     lib.Logger
	repository repository.GeneratorRepository
}

func NewGeneratorService(logger lib.Logger,
	repository repository.GeneratorRepository,
) domains.GeneratorService {
	return GeneratorService{
		logger:     logger,
		repository: repository,
	}
}

func (gs GeneratorService) Store(generator models.GeneratorStoreRequest, userID uint) (models.GeneratorResponse, error) {
	amount := decimal.NewFromFloat(generator.Amount)

	dateFrom, err := time.Parse(constants.DateFormat, generator.DateFrom)
	if err != nil {
		return models.GeneratorResponse{}, err
	}

	var dateTo *sql.NullTime
	if generator.DateTo != nil {
		date, err := time.Parse(constants.DateFormat, *generator.DateTo)
		if err != nil {
			return models.GeneratorResponse{}, err
		}
		dateTo = &sql.NullTime{Time: date, Valid: true}
	}

	model := &models.Generator{
		UserID:            userID,
		Title:             generator.Title,
		Amount:            amount,
		Periodicity:       generator.Periodicity,
		PeriodicityFactor: generator.PeriodicityFactor,
		BudgetFrom:        convertBudgetIDToModel(generator.BudgetFrom),
		BudgetTo:          convertBudgetIDToModel(generator.BudgetTo),
		DateFrom:          dateFrom,
		DateTo:            dateTo,
	}

	if err := gs.repository.Store(model); err != nil {
		return models.GeneratorResponse{}, err
	}

	resp := models.GeneratorResponse{
		ID:                model.ID,
		Title:             model.Title,
		Amount:            model.Amount,
		Periodicity:       model.Periodicity,
		PeriodicityFactor: model.PeriodicityFactor,
		BudgetFrom:        convertBudgetIDFromModel(model.BudgetFrom),
		BudgetTo:          convertBudgetIDFromModel(model.BudgetTo),
		DateFrom:          model.DateFrom.Format(constants.DateFormat),
		DateTo:            convertNullTime(model.DateTo),
	}

	return resp, nil
}

func (gs GeneratorService) List(userID uint) ([]models.GeneratorResponse, error) {
	gens, err := gs.repository.List(userID)
	if err != nil {
		return nil, err
	}

	var resp []models.GeneratorResponse
	for _, v := range gens {
		gen := models.GeneratorResponse{
			ID:                v.ID,
			Title:             v.Title,
			Amount:            v.Amount,
			Periodicity:       v.Periodicity,
			PeriodicityFactor: v.PeriodicityFactor,
			BudgetFrom:        convertBudgetIDFromModel(v.BudgetFrom),
			BudgetTo:          convertBudgetIDFromModel(v.BudgetTo),
			DateFrom:          v.DateFrom.Format(constants.DateFormat),
			DateTo:            convertNullTime(v.DateTo),
		}
		resp = append(resp, gen)
	}
	return resp, nil
}

func (gs GeneratorService) Get(c *gin.Context, userID uint) (models.GeneratorResponse, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.GeneratorResponse{}, err
	}

	gen, err := gs.repository.Get(uint(id), userID)
	if err != nil {
		return models.GeneratorResponse{}, err
	}

	resp := models.GeneratorResponse{
		ID:                gen.ID,
		Title:             gen.Title,
		Amount:            gen.Amount,
		Periodicity:       gen.Periodicity,
		PeriodicityFactor: gen.PeriodicityFactor,
		BudgetFrom:        convertBudgetIDFromModel(gen.BudgetFrom),
		BudgetTo:          convertBudgetIDFromModel(gen.BudgetTo),
		DateFrom:          gen.DateFrom.Format(constants.DateFormat),
		DateTo:            convertNullTime(gen.DateTo),
	}

	return resp, nil
}

func (gs GeneratorService) Update(c *gin.Context, generator models.GeneratorPatchRequest, userID uint) (models.GeneratorResponse, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.GeneratorResponse{}, err
	}

	var dateTo time.Time
	if generator.DateTo != "" {
		dateTo, err = time.Parse(constants.DateFormat, generator.DateTo)
		if err != nil {
			return models.GeneratorResponse{}, err
		}
	}

	var dateFrom time.Time
	if generator.DateFrom != "" {
		dateFrom, err = time.Parse(constants.DateFormat, generator.DateFrom)
		if err != nil {
			return models.GeneratorResponse{}, err
		}
	}

	var amount decimal.Decimal
	if generator.Amount != 0 {
		amount = decimal.NewFromFloat(generator.Amount)
	}

	gen := models.Generator{
		Title:             generator.Title,
		Amount:            amount,
		Periodicity:       generator.Periodicity,
		PeriodicityFactor: generator.PeriodicityFactor,
		BudgetFrom:        convertBudgetIDToModel(generator.BudgetFrom),
		BudgetTo:          convertBudgetIDToModel(generator.BudgetTo),
		DateTo:            &sql.NullTime{Time: dateTo, Valid: true},
		DateFrom:          dateFrom,
	}

	model, err := gs.repository.Update(gen, uint(id), userID)
	if err != nil {
		return models.GeneratorResponse{}, err
	}

	resp := models.GeneratorResponse{
		ID:                model.ID,
		Title:             model.Title,
		Amount:            model.Amount,
		Periodicity:       model.Periodicity,
		PeriodicityFactor: model.PeriodicityFactor,
		BudgetFrom:        convertBudgetIDFromModel(model.BudgetFrom),
		BudgetTo:          convertBudgetIDFromModel(model.BudgetTo),
		DateFrom:          model.DateFrom.Format(constants.DateFormat),
		DateTo:            convertNullTime(model.DateTo),
	}

	return resp, nil
}

func (gs GeneratorService) Delete(c *gin.Context, userID uint) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	err = gs.repository.Delete(uint(id), userID)
	if err != nil {
		return err
	}

	return nil
}

func convertBudgetIDToModel(id *uint) *sql.NullInt64 {
	if id != nil {
		return &sql.NullInt64{
			Int64: int64(*id),
			Valid: true,
		}
	}
	return &sql.NullInt64{}
}

func convertBudgetIDFromModel(budget *sql.NullInt64) *uint {
	if budget == nil {
		return nil
	}
	if budget.Valid {
		var id uint
		id = uint(budget.Int64)
		return &id
	}
	return nil
}

func convertNullTime(time *sql.NullTime) *string {
	if time == nil {
		return nil
	}
	if time.Valid {
		tmp := time.Time.Format(constants.DateFormat)
		return &tmp
	}
	return nil
}
