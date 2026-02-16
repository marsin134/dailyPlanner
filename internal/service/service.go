package service

import (
	"context"
	"dailyPlanner/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type authServiceInterface interface {
	generateAccessToken(user *models.User, session *models.UserSessions) (string, error)
	generateRefreshToken() (string, time.Time, error)
	CheckUserAgentAndIp(sessions []*models.UserSessions, userAgent, ipAddress string) *models.UserSessions
	Register(ctx context.Context, req createUserRequest, ipAddress string) (*models.User, error)
	CreateUserSessionsService(ctx context.Context, user *models.User, ipAddress, userAgent string) (*models.UserSessions, error)
	Login(ctx context.Context, req loginUserRequest, userAgent, ipAddress string) (*models.User, string, *models.UserSessions, error)
	ValidateToken(accessToken string) (*jwt.Token, error)
	GetUserAndSessionFromToken(accessToken string) (*models.User, *models.UserSessions, error)
}

type Service struct {
	AuthService authServiceInterface
}
