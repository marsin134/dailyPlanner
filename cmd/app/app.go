package app

import (
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/database"
	"dailyPlanner/internal/repository"
	"dailyPlanner/internal/service"
	"log"
)

func App(cfg *config.Config) (*database.DB, *repository.Repository, *service.Service) {
	// connection DB
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, nil, nil
	}

	repo := repository.NewRepository(db)
	svc := service.NewService(repo.User, repo.UserSessions, cfg)

	return db, repo, svc
}
