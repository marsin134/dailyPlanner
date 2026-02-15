package repository

import (
	"context"
	"dailyPlanner/internal/database"
	"dailyPlanner/internal/models"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userSessionsRepository struct {
	db *database.DB
}

func NewUserSessionsRepository(db *database.DB) *userSessionsRepository {
	return &userSessionsRepository{db: db}
}

func (s userSessionsRepository) CreateUserSessions(ctx context.Context, session models.UserSessions, refreshToken string) (string, error) {
	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("CreateUserSessions: Error generating bcrypt hash for refresh token: %w", err)
	}

	session.SessionId = uuid.New().String()
	session.RefreshTokenHash = string(refreshTokenHash)
	session.IsActive = true
	session.CreatedAt = time.Now()

	query := `INSERT
		INTO user_sessions (session_id, user_id, refresh_token_hash, expires_at, user_agent, ip_address, is_active, created_at)
		VALUES (:session_id, :user_id, :refresh_token_hash, :expires_at, :user_agent, :ip_address, :is_active, :created_at)`

	_, err = s.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return "", fmt.Errorf("CreateUserSessions: %w", err)
	}

	return session.SessionId, nil
}

func (s userSessionsRepository) GetSessionById(ctx context.Context, sessionId string) (*models.UserSessions, error) {
	query := `SELECT * FROM user_sessions WHERE session_id = $1`

	var session models.UserSessions
	err := s.db.GetContext(ctx, &session, query, sessionId)
	if err != nil {
		return nil, fmt.Errorf("GetUserSessionsById: Error getting session by token hash from db: %w", err)
	}
	return &session, nil
}

func (s userSessionsRepository) GetSessionsByUser(ctx context.Context, userId string) ([]*models.UserSessions, error) {
	query := `SELECT * FROM user_sessions WHERE user_id = $1`

	var sessions []*models.UserSessions
	err := s.db.SelectContext(ctx, &sessions, query, userId)
	if err != nil {
		return nil, fmt.Errorf("GetSessionsByUses: %w", err)
	}
	return sessions, nil
}

func (s userSessionsRepository) UpdateSessionsToken(ctx context.Context, sessionId, newRefreshToken string, expiresAt time.Time) error {
	newRefreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("UpdateSessions: Error generating bcrypt hash for refresh token: %w", err)
	}
	session, err := s.GetSessionById(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("UpdateSessions: Error getting sessions: %w", err)
	}

	session.RefreshTokenHash = string(newRefreshTokenHash)
	session.ExpiresAt = expiresAt

	query := `UPDATE user_sessions 
	SET refresh_token_hash = :refresh_token_hash, expires_at = :expires_at 
                     WHERE session_id = :session_id`

	result, err := s.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("UpdateSessions: Error updating sessions: %w", err)
	}

	if !(CheckUpdate(result)) {
		return fmt.Errorf("UpdateSessions: Error updating sessions")
	}

	return nil
}

func (s userSessionsRepository) Deactivate(ctx context.Context, sessionId string) error {
	query := `UPDATE user_sessions SET is_active = false WHERE session_id = $1`

	result, err := s.db.ExecContext(ctx, query, sessionId)
	if err != nil {
		return fmt.Errorf("deactivate: Error updating sessions: %w", err)
	}

	if !(CheckUpdate(result)) {
		return fmt.Errorf("deactivate: Error updating sessions")
	}
	return nil
}

func (s userSessionsRepository) DeactivateAllExcept(ctx context.Context, userID, currentSessionId string) error {
	query := `UPDATE user_sessions SET is_active = false WHERE user_id = $1 AND session_id != $2`

	result, err := s.db.ExecContext(ctx, query, userID, currentSessionId)
	if err != nil {
		return fmt.Errorf("deactivate: Error updating sessions: %w", err)
	}

	if !(CheckUpdate(result)) {
		return fmt.Errorf("deactivate: Error updating sessions")
	}
	return nil
}

func (s userSessionsRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE is_active = false`

	result, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("DeleteExpired: Error deleting expired sessions: %w", err)
	}

	if !(CheckUpdate(result)) {
		return fmt.Errorf("DeleteExpired: Error deleting expired sessions")
	}
	return nil
}
