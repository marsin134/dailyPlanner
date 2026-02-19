package app

import (
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/handler"
	"dailyPlanner/internal/repository"
	"dailyPlanner/internal/service"
	"fmt"
	"net/http"
)

func InitializationHandlers(repo *repository.Repository, svc *service.Service, cfg *config.Config) http.Handler {
	handlers := handler.Handler{Service: svc, Repo: repo, Cfg: cfg}

	mux := http.NewServeMux()

	mux.HandleFunc("/", helloWorld)

	mux.HandleFunc("/api/auth/register", handlers.Register)
	mux.HandleFunc("/api/auth/login", handlers.LoginHandler)
	mux.HandleFunc("/api/auth/refresh-token", handlers.RefreshToken)
	mux.HandleFunc("/api/auth/logout", handlers.LogoutHandler)
	mux.HandleFunc("/api/auth/logout-all", handlers.LogoutAllExceptHandler)

	mux.HandleFunc("/api/me", handlers.GetMe)
	mux.HandleFunc("/api/me/update-name", handlers.UpdateUserNameHandler)
	mux.HandleFunc("/api/me/update-password", handlers.UpdateUserPasswordHandler)
	mux.HandleFunc("/api/me/delete", handlers.DeleteUser)

	mux.HandleFunc("/api/user/", handlers.GetByUserIDHandler)
	mux.HandleFunc("/api/user/appointment-moderator/", handlers.AppointmentModeratorHandler)

	mux.HandleFunc("/api/event/create", handlers.CreateEventHandler)
	mux.HandleFunc("/api/event/get/", handlers.GetEventByIDHandler)
	mux.HandleFunc("/api/event/get", handlers.GetEventByUserAndDate)
	mux.HandleFunc("/api/event/complete/", handlers.CompleteEvent)
	mux.HandleFunc("/api/event/update", handlers.UpdateEventHandler)
	mux.HandleFunc("api/event/delete/", handlers.DeleteEventByIDHandler)

	return mux
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
