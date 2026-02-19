package middleware

import (
	"context"
	"dailyPlanner/internal/config"
	"dailyPlanner/internal/handler"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

// AuthMiddleware verifies the JWT token and adds user data to the context
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skipping public endpoints
			publicPaths := []string{
				"/api/auth/register",
				"/api/auth/login",
				"/api/auth/refresh-token",
				"/",
			}

			for _, path := range publicPaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Extracting the token from the header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handler.WriteErrorResponse(w, "Authorization is required", http.StatusUnauthorized)
				return
			}

			// Checking the "Bearer <token>" format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handler.WriteErrorResponse(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Parse token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Checking the signature algorithm
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
				}
				return []byte(cfg.Token.JWTSecret), nil
			})

			if err != nil {
				handler.WriteErrorResponse(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				handler.WriteErrorResponse(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Extracting claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				sessionID, ok := claims["session_id"].(string)

				if !ok {
					handler.WriteErrorResponse(w, "Incorrect data in the token", http.StatusUnauthorized)
					return
				}

				// Adding user data to the context
				ctx := r.Context()
				ctx = context.WithValue(ctx, "sessionId", sessionID)
				// Passing the updated context on
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				handler.WriteErrorResponse(w, "Неверные claims токена", http.StatusUnauthorized)
			}
		})
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, URl: %s\nBody: %s\nContext: %s\n\n", r.Method, r.RequestURI, r.Body, r.Context())
		next.ServeHTTP(w, r)
	})
}

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
