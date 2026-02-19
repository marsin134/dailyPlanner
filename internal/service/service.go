package service

import (
	"context"
	"dailyPlanner/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type authServiceInterface interface {
	generateAccessToken(session *models.UserSessions) (string, error)
	generateRefreshToken() (string, time.Time, error)
	CheckUserAgentAndIp(sessions []*models.UserSessions, userAgent, ipAddress string) *models.UserSessions
	Register(ctx context.Context, req CreateUserRequest) (*models.User, error)
	CreateUserSessionsService(ctx context.Context, user *models.User, ipAddress, userAgent string) (string, *models.UserSessions, error)
	Login(ctx context.Context, req LoginUserRequest, userAgent, ipAddress string) (*models.User, string, string, *models.UserSessions, error)
	ValidateToken(accessToken string) (*jwt.Token, error)
	GetSessionFromToken(accessToken string) (*models.UserSessions, error)
	RefreshToken(ctx context.Context, sessionId string) (*models.User, string, string, error)
}

type Service struct {
	AuthService authServiceInterface
}
