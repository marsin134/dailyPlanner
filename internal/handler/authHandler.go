package handler

import (
	"dailyPlanner/internal/service"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

func getIP(r *http.Request) string {
	// Check proxy headers (if behind nginx/cloudflare)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// An alternative header from some proxies
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// If there is no proxy, we take RemoteAddr.
	ip = r.RemoteAddr
	if host, _, err := net.SplitHostPort(ip); err == nil {
		return host
	}

	return ip
}

func getUserAgent(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return "unknown"
	}
	return userAgent
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.CheckHandlerStruct(w)
	if err != nil {
		return
	}

	var req RegisterRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// email verification
	patternEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(patternEmail, req.Email)
	if err != nil || !matched {
		WriteErrorResponse(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// password verification
	if utf8.RuneCountInString(req.Password) < 6 {
		WriteErrorResponse(w, "The password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	if err = h.Validate.Struct(req); err != nil {
		WriteErrorResponse(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	registerRequest := service.CreateUserRequest{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err = h.Service.AuthService.Register(r.Context(), registerRequest)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginRequest := service.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	userAgent := getUserAgent(r)
	userIp := getIP(r)
	user, accessToken, refreshToken, session, err := h.Service.AuthService.Login(r.Context(), loginRequest, userAgent, userIp)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			UserID:   user.UserId,
			UserName: user.UserName,
			Email:    user.Email,
		},
		Session: SessionResponse{
			SessionID: session.SessionId,
			UserID:    session.UserId,
			ExpiresAt: session.ExpiresAt,
			UserAgent: session.UserAgent,
			IpAddress: session.IpAddress,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.CheckHandlerStruct(w) != nil {
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		WriteErrorResponse(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	loginRequest := service.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	userAgent := getUserAgent(r)
	userIp := getIP(r)
	user, accessToken, refreshToken, session, err := h.Service.AuthService.Login(r.Context(), loginRequest, userAgent, userIp)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			UserID:   user.UserId,
			UserName: user.UserName,
			Email:    user.Email,
		},
		Session: SessionResponse{
			SessionID: session.SessionId,
			ExpiresAt: session.ExpiresAt,
			UserAgent: session.UserAgent,
			IpAddress: session.IpAddress,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	if h.CheckHandlerStruct(w) != nil {
		return
	}

	var req RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.SessionId == "" {
		WriteErrorResponse(w, "Session id is required", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		WriteErrorResponse(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	session, err := h.Repo.UserSessions.GetSessionById(r.Context(), req.SessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(session.RefreshTokenHash), []byte(req.RefreshToken))
	if err != nil {
		WriteErrorResponse(w, "Refresh token doesn't match", http.StatusUnauthorized)
		return
	}

	user, accessToken, RefreshToken, err := h.Service.AuthService.RefreshToken(r.Context(), req.SessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err = h.Repo.UserSessions.GetSessionById(r.Context(), req.SessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: RefreshToken,
		User: UserResponse{
			UserID:   user.UserId,
			UserName: user.UserName,
			Email:    user.Email,
		},
		Session: SessionResponse{
			SessionID: req.SessionId,
			ExpiresAt: session.ExpiresAt,
			UserAgent: session.UserAgent,
			IpAddress: session.IpAddress,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.CheckHandlerStruct(w) != nil {
		return
	}

	sessionId, ok := r.Context().Value("sessionId").(string)
	if !ok {
		WriteErrorResponse(w, "Authorization is required", http.StatusBadRequest)
		return
	}

	err := h.Repo.UserSessions.Deactivate(r.Context(), sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "Successful exit from the device"})
}

func (h *Handler) LogoutAllExceptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.CheckHandlerStruct(w) != nil {
		return
	}

	sessionId, ok := r.Context().Value("sessionId").(string)
	if !ok {
		WriteErrorResponse(w, "Authorization is required", http.StatusBadRequest)
		return
	}

	session, err := h.Repo.UserSessions.GetSessionById(r.Context(), sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Repo.UserSessions.DeactivateAllExcept(r.Context(), session.UserId, sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "Successful exit from the devices"})
}
