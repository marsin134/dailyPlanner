package service

import (
	"context"
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/models"
	"dailyPlanner/internal/repository"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type authService struct {
	userRepo     repository.UserRepository
	sessionsRepo repository.UserSessionsRepository
	cfg          *config.Config
}

func NewAuthService(userRepo repository.UserRepository, sessionsRepo repository.UserSessionsRepository, cfg *config.Config) *authService {
	return &authService{userRepo: userRepo, sessionsRepo: sessionsRepo, cfg: cfg}
}

type createUserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (svc *authService) generateAccessToken(user *models.User, session *models.UserSessions) (string, error) {
	claims := jwt.MapClaims{
		"user_name":  user.UserName,
		"user_id":    user.UserId,
		"email":      user.Email,
		"role":       user.Role,
		"session_id": session.SessionId,
		"exp":        time.Now().Add(svc.cfg.Token.AccessTokenDuration).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(svc.cfg.Token.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("ошибка подписи токена: %w", err)
	}

	return tokenString, nil
}

func (svc *authService) generateRefreshToken() (string, time.Time, error) {
	refreshToken := uuid.New().String()

	expiryTime := time.Now().Add(svc.cfg.Token.RefreshTokenDuration)

	return refreshToken, expiryTime, nil
}

func (svc *authService) CheckUserAgentAndIp(sessions []*models.UserSessions, userAgent, ipAddress string) *models.UserSessions {
	for _, session := range sessions {
		if session.UserAgent == userAgent && session.IpAddress == ipAddress {
			return session
		}
	}
	return nil
}

func (svc *authService) Register(ctx context.Context, req createUserRequest, ipAddress string) (*models.User, error) {
	// get user by email
	existingUser, err := svc.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("the user with email %s already exists", req.Email)
	}

	user := &models.User{
		UserName: req.UserName,
		Email:    req.Email,
	}

	err = svc.userRepo.CreateUser(ctx, user, req.Password)
	if err != nil {
		return nil, err
	}

	user, err = svc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *authService) CreateUserSessionsService(ctx context.Context, user *models.User, ipAddress, userAgent string) (*models.UserSessions, error) {
	// creating a refresh token
	refreshToken, expiresAt, err := svc.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("error when creating a user: %w", err)
	}

	sessionUser := models.UserSessions{
		UserId:    user.UserId,
		ExpiresAt: expiresAt,
		UserAgent: userAgent,
		IpAddress: ipAddress}

	sessionId, err := svc.sessionsRepo.CreateUserSessions(ctx, sessionUser, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("error when creating a user: %w", err)
	}

	sessionAgent, err := svc.sessionsRepo.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("error when creating a user: %w", err)
	}

	return sessionAgent, nil
}

func (svc *authService) Login(ctx context.Context, req loginUserRequest, userAgent, ipAddress string) (*models.User, string, *models.UserSessions, error) {
	// checking for password compliance
	user, err := svc.userRepo.VerifyPassword(ctx, req.Email, req.Password)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error when logging in: %w", err)
	}

	// getting all user sessions
	sessions, err := svc.sessionsRepo.GetSessionsByUser(ctx, user.UserId)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error when logging in: %w", err)
	}

	if len(sessions) == 0 {
		session, err := svc.CreateUserSessionsService(ctx, user, ipAddress, userAgent)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error when logging in: %w", err)
		}
		return user, "", session, nil
	}

	session := svc.CheckUserAgentAndIp(sessions, userAgent, ipAddress)
	if session == nil {
		session, err = svc.CreateUserSessionsService(ctx, user, ipAddress, userAgent)
	} else {
		// updating a refresh token
		refreshToken, expiresAt, err := svc.generateRefreshToken()
		if err != nil {
			return nil, "", nil, fmt.Errorf("error when logging in: %w", err)
		}
		err = svc.sessionsRepo.UpdateSessionsToken(ctx, session.SessionId, refreshToken, expiresAt)
	}

	if err != nil {
		return nil, "", nil, fmt.Errorf("error when logging in: %w", err)
	}

	accessToken, err := svc.generateAccessToken(user, session)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error when logging in: %w", err)
	}

	return user, accessToken, session, nil
}

func (svc *authService) ValidateToken(accessToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(svc.cfg.Token.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error when validating token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (svc *authService) GetUserAndSessionFromToken(accessToken string) (*models.User, *models.UserSessions, error) {
	token, err := svc.ValidateToken(accessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("error when validating token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("error when validating token")
	}

	user := &models.User{
		UserName: claims["user_name"].(string),
		Email:    claims["email"].(string),
	}

	session := &models.UserSessions{
		SessionId: claims["session_id"].(string),
		UserId:    claims["user_id"].(string),
	}

	return user, session, nil
}
