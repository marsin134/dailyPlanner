package repository

import (
	"context"
	"dailyPlanner/internal/models"
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
