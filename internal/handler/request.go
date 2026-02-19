package handler

// Auth request

type RegisterRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email" Validate:"required,email"`
	Password string `json:"password" Validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" Validate:"required,email"`
	Password string `json:"password" Validate:"required,min=6"`
}

type RefreshTokenRequest struct {
	SessionId    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
}

// User request

type UpdateUserNameRequest struct {
	NewUserName string `json:"user_name"`
}

type UpdateUserPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// Event request

type CreateEventRequest struct {
	Title     string `json:"title"`
	DateEvent string `json:"date_event"`
	Color     string `json:"color"`
}

type GetEventsRequestForDate struct {
	DateEvent string `json:"date_event"`
}

type UpdateEventRequest struct {
	EventId string `json:"event_id"`
	Title   string `json:"title"`
	Color   string `json:"color"`
}
