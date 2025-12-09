package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ServiceName string `json:"service_name"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens for internal service communication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "missing authorization header",
				"msg":   "Authorization header with Bearer token is required",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "invalid authorization format",
				"msg":   "Authorization header must be in format: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			secretEnv := os.Getenv("JWT_SECRET")
			return []byte(secretEnv), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "invalid token",
				"msg":   err.Error(),
			})
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "token validation failed",
				"msg":   "token is not valid",
			})
			return
		}

		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "token expired",
				"msg":   "please obtain a new token",
			})
			return
		}

		// Token is valid, proceed to next handler
		next.ServeHTTP(w, r)
	})
}
