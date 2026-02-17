package handler

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

type LogoutRequest struct {
	SessionId string `json:"session_id"`
}

type LogoutAllExcept struct {
	ActiveSessionId string `json:"active_session_id"`
	UserId          string `json:"user_id"`
}
