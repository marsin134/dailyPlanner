package handler

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse - standard response with an error
type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{message})
}
