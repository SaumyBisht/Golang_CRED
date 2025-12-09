package utils

import (
	"deployment-service/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
	ServiceName string `json:"service_name"`
	jwt.RegisteredClaims
}

// GenerateServiceToken creates a JWT token for service-to-service communication
func GenerateServiceToken() (string, error) {
	// Token valid for 1 hour
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		ServiceName: config.ServiceName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.ServiceName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
