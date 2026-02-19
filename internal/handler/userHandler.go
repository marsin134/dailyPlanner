package handler

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionId, ok := r.Context().Value("sessionId").(string)
	if !ok {
		WriteErrorResponse(w, "Session Id not found in request context", http.StatusUnauthorized)
		return
	}

	session, err := h.Repo.UserSessions.GetSessionById(r.Context(), sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.Repo.User.GetUserById(r.Context(), session.UserId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := UserResponse{
		UserID:   user.UserId,
		UserName: user.UserName,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionId, ok := r.Context().Value("sessionId").(string)
	if !ok {
		WriteErrorResponse(w, "Session Id not found in request context", http.StatusUnauthorized)
		return
	}

	session, err := h.Repo.UserSessions.GetSessionById(r.Context(), sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.Repo.User.GetUserById(r.Context(), session.UserId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Role != "Admin" {
		WriteErrorResponse(w, "You do not have permission to view this resource", http.StatusForbidden)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	userID := pathParts[len(pathParts)-1]

	user, err = h.Repo.User.GetUserById(r.Context(), userID)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
	}

	response := UserResponse{
		UserID:   user.UserId,
		UserName: user.UserName,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateUserNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	var req UpdateUserNameRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		WriteErrorResponse(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	sessionId, ok := r.Context().Value("sessionId").(string)
	if !ok {
		WriteErrorResponse(w, "Session Id not found in request context", http.StatusUnauthorized)
		return
	}

	session, err := h.Repo.UserSessions.GetSessionById(r.Context(), sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.Repo.User.GetUserById(r.Context(), session.UserId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.User.UpdateUsername(r.Context(), user.Email, req.NewUserName)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "The user has been updated"})
}

func (h *Handler) UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	var req UpdateUserPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		WriteErrorResponse(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	sessionId, ok := r.Context().Value("sessionId").(string)
	if !ok {
		WriteErrorResponse(w, "Session Id not found in request context", http.StatusUnauthorized)
		return
	}

	session, err := h.Repo.UserSessions.GetSessionById(r.Context(), sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.Repo.User.GetUserById(r.Context(), session.UserId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		WriteErrorResponse(w, "Passwords don't match", http.StatusBadRequest)
	}

	err = h.Repo.User.UpdatePassword(r.Context(), user.Email, req.OldPassword, req.NewPassword)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "The user has been updated"})
}

func (h *Handler) AppointmentModeratorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	userId := pathParts[len(pathParts)-1]

	sessionId, ok := r.Context().Value("sessionId").(string)
	if !ok {
		WriteErrorResponse(w, "Session Id not found in request context", http.StatusUnauthorized)
		return
	}

	session, err := h.Repo.UserSessions.GetSessionById(r.Context(), sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.Repo.User.GetUserById(r.Context(), session.UserId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Role != "Admin" {
		WriteErrorResponse(w, "You do not have permission to moderate this resource", http.StatusForbidden)
		return
	}

	err = h.Repo.User.AppointmentModerator(r.Context(), userId, "Admin")
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "The user has been applied"})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	sessionId, ok := r.Context().Value("sessionId").(string)
	if !ok {
		WriteErrorResponse(w, "Session Id not found in request context", http.StatusUnauthorized)
		return
	}

	session, err := h.Repo.UserSessions.GetSessionById(r.Context(), sessionId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.Repo.User.DeleteUser(r.Context(), session.UserId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "The user has been deleted"})
}
