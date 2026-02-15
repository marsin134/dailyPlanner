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
		log.Printf("Не удалось загрузить .env: %v", err)
	}

	cfg := config.LoadConfig()

	if cfg.Token.JWTSecret == "" {
		log.Fatal("JWT_SECRET_KEY не установлен в .env файле")
	}

	db := app.App(&cfg)
	defer database.MethodsDB.Close(db)

	//ctx := context.Context(context.Background())
	//
	//user := models.User{
	//	UserId:   "1",
	//	UserName: "Oleg",
	//	Email:    "oleg@gmail.com",
	//	Role:     "User",
	//}
	//r := repository.NewUserRepository(db)
	//rs := repository.NewUserSessionsRepository(db)
	//err := r.CreateUser(ctx, &user, "123")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//getUser, err := r.GetUserByEmail(ctx, "oleg@gmail.com")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//sessions, err := rs.GetSessionsByUser(ctx, getUser.UserId)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//println(sessions, len(sessions))

	// Starting the server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	fmt.Printf("Сервер запущен на %s\n", addr)

	http.HandleFunc("/", helloWorld)
	http.ListenAndServe(addr, nil)
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
