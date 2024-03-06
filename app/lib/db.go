package lib

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"finapp/models"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(env Env, logger Logger) Database {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		env.DBUser,
		env.DBPassword,
		env.DBName,
		env.DBHost,
		env.DBPort,
		env.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.GetGormLogger(),
	})
	if err != nil {
		logger.Info("Connect to sqlite...")
		dsn := "test.db"
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.GetGormLogger(),
		})
		if err != nil {
			logger.Panic("Can't open DB: ", err.Error())
		}
	}
	logger.Info("Connected to database")

	if err := db.AutoMigrate(&models.User{}, models.Trx{}, models.Budget{}, models.Goal{}); err != nil {
		logger.Panic("Can't migrate database: ", err.Error())
	}
	logger.Info("Migrated database")

	return Database{
		db,
	}
}
