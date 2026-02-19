package main

import (
	"dailyPlanner/cmd/app"
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/database"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Uploading .env
	if err := config.LoadEnvFile(".env"); err != nil {
		log.Printf("Failed to load .env: %v", err)
	}

	cfg := config.LoadConfig()

	if cfg.Token.JWTSecret == "" {
		log.Fatal("JWT_SECRET_KEY is not set in the .env file")
	}

	db, repo, svc := app.App(&cfg)
	defer database.MethodsDB.Close(db)

	handlerChain := app.InitializationHandlers(repo, svc, &cfg)

	// Starting the server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	fmt.Printf("The server is running on %s\n", addr)

	if err := http.ListenAndServe(addr, handlerChain); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
