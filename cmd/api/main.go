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
	"time"
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
	//checkEvents(db, ctx)
	checkUserSessions(db, ctx)
}

func checkUserSessions(db *database.DB, ctx context.Context) {
	s := repository.NewUserSessionsRepository(db)

	date := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	//userSession := models.UserSessions{UserId: "8c522305-7aaf-434c-9a5c-58c1166a58be",
	//	ExpiresAt: date,
	//	IpAddress: ""}
	//
	//err := s.CreateUserSessions(ctx, userSession, "hello")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//println("Good Create")
	//
	//err = s.CreateUserSessions(ctx, userSession, "helloSecond")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//println("Good Create Second")

	session, err := s.GetSessionById(ctx, "905d5df7-4f8b-4fe6-8d36-a013989aa9a3")
	if err != nil {
		log.Fatal(err)
	}
	println("Good GetSessionByTokenHash")

	secondSessions, err := s.GetSessionsByUser(ctx, "8c522305-7aaf-434c-9a5c-58c1166a58be")
	if err != nil {
		log.Fatal(err)
	}

	println("Test sessions")
	for _, secondSession := range secondSessions {
		println(secondSession.UserId)
	}

	println("Good GetSessionsByUser")

	err = s.UpdateSessionsToken(ctx, "f0adc2ed-f7c7-4620-8006-2c63018c480e", "helloSecond2", date)
	if err != nil {
		log.Fatal(err)
	}
	println("Good UpdateSessionsToken")

	err = s.DeactivateAllExcept(ctx, "8c522305-7aaf-434c-9a5c-58c1166a58be", "905d5df7-4f8b-4fe6-8d36-a013989aa9a3")
	if err != nil {
		log.Fatal(err)
	}

	secondSessions, err = s.GetSessionsByUser(ctx, "8c522305-7aaf-434c-9a5c-58c1166a58be")
	if err != nil {
		log.Fatal(err)
	}

	println("Test sessions")
	for _, secondSession := range secondSessions {
		println(secondSession.IsActive)
	}
	println("Good DeactivateAllExcept")

	err = s.Deactivate(ctx, "905d5df7-4f8b-4fe6-8d36-a013989aa9a3")
	if err != nil {
		log.Fatal(err)
	}
	session, _ = s.GetSessionById(ctx, "905d5df7-4f8b-4fe6-8d36-a013989aa9a3")
	println(session.IsActive)
	println("Good Deactivate")

	err = s.DeleteExpired(ctx)
	if err != nil {
		log.Fatal(err)
	}
	println("Good DeleteExpired")
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
