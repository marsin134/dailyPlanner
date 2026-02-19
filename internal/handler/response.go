package handler

import (
	"dailyPlanner/internal/models"
	"time"
)

// Auth Response

type UserResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type SessionResponse struct {
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	UserAgent string    `json:"user_agent"`
	IpAddress string    `json:"ip_address"`
}

type AuthResponse struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	User         UserResponse    `json:"user"`
	Session      SessionResponse `json:"session"`
}

// Event response

type EventResponse struct {
	EventID   string `json:"event_id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
	Color     string `json:"color"`
}

type EventsResponse struct {
	Events []models.Event `json:"events"`
}
