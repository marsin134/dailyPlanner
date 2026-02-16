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
	SessionId string `json:"refresh_token"`
}
