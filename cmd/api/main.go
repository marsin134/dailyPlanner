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

	//checkUser(db, ctx)
	checkEvents(db, ctx)

}

func checkEvents(db *database.DB, ctx context.Context) {
	r := repository.NewUserRepository(db)
	user, _ := r.GetUserById(ctx, "de891e94-f386-43fd-8698-d0721efb3af3")
	event := models.Event{
		TitleEvent: "помыть посуду",
		DateEvent:  "2026-02-09",
	}

	vr := repository.NewEventRepository(db)

	err := vr.CreateEvent(ctx, user.UserId, &event)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Create event")

	firstEvent, err := vr.GetEventsByUserIdAndDate(ctx, user.UserId, event.DateEvent)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(firstEvent[0])
	fmt.Println("Good get event by user.id and date")
	fmt.Println(firstEvent[0].TitleEvent)

	secondEvent, err := vr.GetEventById(ctx, firstEvent[0].EventId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(secondEvent)
	fmt.Println("Good get event by user.id and date")

	err = vr.UpdateEvent(ctx, secondEvent.EventId, "помыть посуды и убрать со стола", "green")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good update event")

	err = vr.CompleteEvent(ctx, secondEvent.EventId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good complete event")

	err = vr.DeleteEvent(ctx, secondEvent.EventId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good delete event")
}

func checkUser(db *database.DB, ctx context.Context) {
	user := models.User{
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
