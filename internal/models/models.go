package models

import "time"

type User struct {
	UserId       string `json:"user_id" db:"user_id"`
	UserName     string `json:"user_name" db:"user_name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password" db:"password_hash"`
	Role         string `json:"role" db:"role"`
}

type UserSessions struct {
	SessionId        string    `json:"session_id" db:"session_id"`
	UserId           string    `json:"user_id" db:"user_id"`
	RefreshTokenHash string    `json:"refresh_token_hash" db:"refresh_token_hash"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	IpAddress        string    `json:"ip_address" db:"ip_address"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type Event struct {
	EventId    string `json:"event_id" db:"event_id"`
	UserId     string `json:"user_id" db:"user_id"`
	TitleEvent string `json:"title_event" db:"title_event"`
	DateEvent  string `json:"date_event" db:"date_event"`
	Completed  bool   `json:"completed" db:"completed"`
	Color      string `json:"color" db:"color"`
}
