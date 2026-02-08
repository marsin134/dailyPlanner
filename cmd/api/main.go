package main

import (
	"context"
	"dailyPlanner/cmd/app"
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/database"
	"dailyPlanner/internal/models"
	"dailyPlanner/internal/repository"
	"fmt"
	"log"
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

	ctx := context.Context(context.Background())

	user := models.User{
		UserId:   "1",
		UserName: "Oleg",
		Email:    "oleg@gmail.com",
		Role:     "User",
	}
	r := repository.NewUserRepository(db)
	err := r.CreateUser(ctx, &user, "123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Create user")

	newUser, err := r.GetUserByEmail(ctx, user.Email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(newUser)
	secondUser, err := r.GetUserById(ctx, user.UserId)
	fmt.Println(secondUser)

	userPassword, err := r.VerifyPassword(ctx, user.Email, "123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userPassword)
	_, err = r.VerifyPassword(ctx, user.Email, "notGoodPassword")
	if err != nil {
		fmt.Println("Good")
	}

	err = r.UpdateUsername(ctx, user.Email, "megaOleg", "123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.UserName)
	fmt.Println("Good update user")

	err = r.UpdatePassword(ctx, user.Email, "123", "1234")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good update user password")

	err = r.AppointmentModerator(ctx, user.Email, "Admin", "1234")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good update user role")

	err = r.DeleteUser(ctx, newUser.UserId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good delete user")
}
