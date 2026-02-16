package handler

import (
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/repository"
	"dailyPlanner/internal/service"
	"gitlab.com/golang-library/go-validator"
)

type Handler struct {
	Service  service.Service
	Repo     repository.Repository
	Cfg      *config.Config
	Validate *validator.Validate
}

func NewHandler(service service.Service, repo repository.Repository, cfg *config.Config) *Handler {
	return &Handler{
		Service:  service,
		Repo:     repo,
		Cfg:      cfg,
		Validate: validator.New(),
	}
}
