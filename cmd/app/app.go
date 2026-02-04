package app

import (
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/database"
	"log"
)

func App(cfg *config.Config) *database.DB {
	// connection DB
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
		return nil
	}

	return db
}
