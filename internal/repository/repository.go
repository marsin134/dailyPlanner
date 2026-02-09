package repository

import (
	"context"
	"dailyPlanner/internal/models"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User, password string) error
	GetUserById(ctx context.Context, userId string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	VerifyPassword(ctx context.Context, email string, password string) (*models.User, error)
	UpdateUsername(ctx context.Context, email, newUserName, password string) error
	UpdatePassword(ctx context.Context, email, password, newPassword string) error
	AppointmentModerator(ctx context.Context, email, role, password string) error
	DeleteUser(ctx context.Context, userId string) error
}

type EventRepository interface {
	CreateEvent(ctx context.Context, userId string, event *models.Event) error
	GetEventById(ctx context.Context, eventId string) (*models.Event, error)
	GetEventsByUserIdAndDate(ctx context.Context, userId string, date time.Time) ([]*models.Event, error)
	UpdateEvent(ctx context.Context, eventId, newTitle, color string) error
	CompleteEvent(ctx context.Context, eventId string) error
	DeleteEvent(ctx context.Context, eventId string) error
}
