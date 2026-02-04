package main

import (
	"dailyPlanner/cmd/app"
	"dailyPlanner/internal/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Uploading .env
	if err := config.LoadEnvFile(".env"); err != nil {
		log.Printf("Не удалось загрузить .env: %v", err)
	}

	cfg := config.LoadConfig()

	if cfg.Token.JWTSecret == "" {
		log.Fatal("JWT_SECRET_KEY не установлен в .env файле")
	}

	app.App(&cfg)

	// Starting the server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	fmt.Printf("Сервер запущен на %s\n", addr)

	http.HandleFunc("/", helloWorld)
	http.ListenAndServe(addr, nil)
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
