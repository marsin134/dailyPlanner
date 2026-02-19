package handler

import (
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/repository"
	"dailyPlanner/internal/service"
	"fmt"
	"gitlab.com/golang-library/go-validator"
	"net/http"
)

type Handler struct {
	Service  *service.Service
	Repo     *repository.Repository
	Cfg      *config.Config
	Validate *validator.Validate
}

func NewHandler(service *service.Service, repo *repository.Repository, cfg *config.Config) *Handler {
	return &Handler{
		Service:  service,
		Repo:     repo,
		Cfg:      cfg,
		Validate: validator.New(),
	}
}

func (h *Handler) CheckHandlerStruct(w http.ResponseWriter) error {
	if h == nil {
		WriteErrorResponse(w, "Handler is nil", http.StatusNotImplemented)
		return fmt.Errorf("Handler is nil. ")
	}

	if h.Validate == nil {
		WriteErrorResponse(w, "Validate is nil", http.StatusNotImplemented)
		return fmt.Errorf("Validate is nil. ")
	}

	return nil
}

type MessageResponse struct {
	Message string `json:"message"`
}
