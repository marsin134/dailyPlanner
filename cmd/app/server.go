package app

import (
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/handler"
	"dailyPlanner/internal/middleware"
	"dailyPlanner/internal/repository"
	"dailyPlanner/internal/service"
	"fmt"
	"net/http"
)

func InitializationHandlers(repo *repository.Repository, svc *service.Service, cfg *config.Config) http.Handler {
	handlers := handler.NewHandler(svc, repo, cfg)

	mux := http.NewServeMux()

	mux.HandleFunc("/", HomeHandler)

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
	mux.HandleFunc("/api/event/delete/", handlers.DeleteEventByIDHandler)

	handlerChain := middleware.Chain(
		mux,
		middleware.CORSMiddleware,
		middleware.AuthMiddleware(cfg),
		middleware.LoggingMiddleware)

	return handlerChain
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "/, HomeHandler\n\n")

	fmt.Fprintf(w, "/api/auth/register, handlers.Register\n")
	fmt.Fprintf(w, "/api/auth/login, handlers.LoginHandler\n")
	fmt.Fprintf(w, "/api/auth/refresh-token, handlers.RefreshToken\n")
	fmt.Fprintf(w, "/api/auth/logout, handlers.LogoutHandler\n")
	fmt.Fprintf(w, "/api/auth/logout-all, handlers.LogoutAllExceptHandler\n\n")

	fmt.Fprintf(w, "/api/me, handlers.GetMe\n")
	fmt.Fprintf(w, "/api/me/update-name, handlers.UpdateUserNameHandler\n")
	fmt.Fprintf(w, "/api/me/update-password, handlers.UpdateUserPasswordHandler\n")
	fmt.Fprintf(w, "/api/me/delete/, handlers.DeleteEventByIDHandler\n\n")

	fmt.Fprintf(w, "/api/user/, handlers.GetByUserIDHandler\n")
	fmt.Fprintf(w, "/api/user/appointment-moderator/, handlers.AppointmentModeratorHandler\n\n")

	fmt.Fprintf(w, "/api/event/create, handlers.CreateEventHandler\n")
	fmt.Fprintf(w, "/api/event/get/, handlers.GetEventByIDHandler\n")
	fmt.Fprintf(w, "/api/event/get, handlers.GetEventByUserAndDate\n")
	fmt.Fprintf(w, "/api/event/complete/, handlers.CompleteEvent\n")
	fmt.Fprintf(w, "/api/event/update, handlers.UpdateEventHandler\n")
	fmt.Fprintf(w, "/api/event/delete/, handlers.DeleteEventByIDHandler\n")
}
