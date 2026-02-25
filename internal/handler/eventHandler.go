package handler

import (
	"dailyPlanner/internal/models"
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	var req CreateEventRequest

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

	event := &models.Event{
		TitleEvent: req.Title,
		DateEvent:  req.DateEvent,
		Color:      req.Color,
	}

	event, err = h.Repo.Event.CreateEvent(r.Context(), user.UserId, event)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := EventResponse{
		EventID:   event.EventId,
		UserID:    event.UserId,
		Title:     event.TitleEvent,
		Date:      event.DateEvent,
		Completed: event.Completed,
		Color:     event.Color,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	eventID := pathParts[len(pathParts)-1]

	event, err := h.Repo.Event.GetEventById(r.Context(), eventID)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
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

	if user.Role != "Admin" && user.UserId != event.UserId {
		WriteErrorResponse(w, "You do not have permission to view this resource", http.StatusForbidden)
		return
	}

	response := EventResponse{
		EventID:   event.EventId,
		UserID:    event.UserId,
		Title:     event.TitleEvent,
		Date:      event.DateEvent,
		Completed: event.Completed,
		Color:     event.Color,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetEventByUserAndDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	user, err := h.Repo.User.GetUserById(r.Context(), session.UserId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req GetEventsRequestForDate

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err = h.Validate.Struct(req); err != nil {
		WriteErrorResponse(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	events, err := h.Repo.Event.GetEventsByUserIdAndDate(r.Context(), user.UserId, req.DateEvent)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := EventsResponse{
		Events: events,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		WriteErrorResponse(w, "Incorrect data", http.StatusBadRequest)
		return
	}

	err := h.Repo.Event.UpdateEvent(r.Context(), req.EventId, req.Title, req.Color)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "Event successfully updated"})
}

func (h *Handler) CompleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	eventID := pathParts[len(pathParts)-1]

	event, err := h.Repo.Event.GetEventById(r.Context(), eventID)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
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

	if user.UserId != event.UserId {
		WriteErrorResponse(w, "You do not have permission to view this resource", http.StatusForbidden)
		return
	}

	err = h.Repo.Event.CompleteEvent(r.Context(), event.EventId)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "Event successfully completed"})
}

func (h *Handler) DeleteEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.CheckHandlerStruct(w); err != nil {
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	eventID := pathParts[len(pathParts)-1]

	event, err := h.Repo.Event.GetEventById(r.Context(), eventID)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
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

	if user.Role != "Admin" && user.UserId != event.UserId {
		WriteErrorResponse(w, "You do not have permission to view this resource", http.StatusForbidden)
		return
	}

	err = h.Repo.Event.DeleteEvent(r.Context(), eventID)
	if err != nil {
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{Message: "Event successfully deleted"})
}
